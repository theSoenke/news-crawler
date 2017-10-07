package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/thesoenke/news-crawler/scraper"
)

var cmdExtract = &cobra.Command{
	Use:   "extract",
	Short: "Extract contents from HTML files in input dir",
	RunE: func(cmd *cobra.Command, args []string) error {
		inputDir := args[0]
		outputDir := args[1]

		stat, err := os.Stat(inputDir)
		if err != nil {
			return err
		}

		if !stat.IsDir() {
			return fmt.Errorf("Input %s needs to be a directory", inputDir)
		}

		err = extractContent(inputDir, outputDir)
		return err
	},
}

func init() {
	cmdScrape.Args = cobra.ExactArgs(2)
	RootCmd.AddCommand(cmdExtract)
}

func extractContent(inputDir string, outputDir string) error {
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		filepath := path.Join(inputDir, f.Name())
		file, err := ioutil.ReadFile(filepath)
		if err != nil {
			return err
		}

		content, err := scraper.ExtractWithGoOse("", string(file))
		if err != nil {
			return err
		}

		filepath = path.Join(outputDir, f.Name())
		err = ioutil.WriteFile(filepath, []byte(content), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
