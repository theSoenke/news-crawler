package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var timezone string
var lang string
var verbose bool
var logsDir string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "news-crawler",
	Short: "News article scraper",
	Long:  `Scraper to extract the content of daily news articles`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func writeLog(filename string, log string) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return err
	}

	_, err = f.Write([]byte(log))
	return err
}
