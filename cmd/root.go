package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var rootCmd = &cobra.Command{
	Use:   "raindrop",
	Short: "CLI client for Raindrop.io",
	Long:  `Raindrop.io is tool to manage bookmarks. It lets you save your bookmarks into folders and those bookmarks are then available from anywhere.`,
	Run:   func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// ~/.raindrop/config.yml
type Config struct {
	TestToken    string `yaml:"testToken"`
	RefreshToken string `yaml:"refresh_token"`
	ClientID     string `yaml:"clientID"`
	ClientSecret string `yaml:"clientSecret"`
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func readFile(filename string) *os.File {
	file, err := os.Open(filename)
	if err != nil {
		processError(err)
		var temp *os.File
		return temp
	}
	return file
}

func parseYaml(cfg *Config, file *os.File) bool {
	decoder := yaml.NewDecoder(file)
	err := decoder.Decode(cfg)
	if err != nil {
		processError(err)
		return false
	}
	return true
}

func getYaml(filename string) Config {
	var cfg Config
	yamlFile := readFile(filename)
	parsed := parseYaml(&cfg, yamlFile)
	if !parsed {
		log.Fatal("Failed to read file.")
	}
	return cfg
}

func pathToConfig() string {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return user.HomeDir + "/.raindrop/config.yml"
}

// ACCESS METHODS
func getAccessToken() string {
	var cfg = getYaml(pathToConfig())
	return cfg.TestToken
}

func getRefreshToken() string {
	var cfg = getYaml(pathToConfig())
	return cfg.RefreshToken
}

func getClientID() string {
	var cfg = getYaml(pathToConfig())
	return cfg.ClientID
}

func getClientSecret() string {
	var cfg = getYaml(pathToConfig())
	return cfg.ClientSecret
}

// GET METHODS
func makeGetRequest(url string) io.ReadCloser {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	token := getAccessToken()
	req.Header.Set("Authorization", "Bearer "+token)
	res, _ := client.Do(req)
	return res.Body
}

// POST METHODS
func makePostRequest(url string, data *bytes.Buffer) int {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, data)
	token := getAccessToken()
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		processError(err)
	}
	return response.StatusCode
}
