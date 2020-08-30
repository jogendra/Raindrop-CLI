package cmd

import (
	"fmt"
	"os"

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
		// TODO: Process error
		var temp *os.File
		return temp
	}
	return file
}

func parseYaml(cfg *Config, file *os.File) bool {
	decoder := yaml.NewDecoder(file)
	err := decoder.Decode(cfg)
	if err != nil {
		// TODO: process error
		return false
	}
	return true
}
