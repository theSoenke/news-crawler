package cmd

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/thesoenke/news-crawler/scraper"
)

var scrapeOutDir string
var scrapeVerbose bool

var cmdScrape = &cobra.Command{
	Use:   "scrape",
	Short: "Scrape all provided articles",
	RunE: func(cmd *cobra.Command, args []string) error {
		itemsInputFile := args[0]

		location, err := time.LoadLocation(timezone)
		if err != nil {
			return err
		}

		path, err := getFeedsFilePath(itemsInputFile, location)
		if err != nil {
			return err
		}

		contentScraper, err := scraper.New(path)
		if err != nil {
			return err
		}

		elasticClient, err := scraper.NewElasticClient()
		if err != nil {
			return err
		}

		start := time.Now()
		err = contentScraper.Scrape(scrapeOutDir, location, elasticClient, scrapeVerbose)
		if err != nil {
			return err
		}

		log.Printf("Articles: %d successful, %d failures in %s from %s", contentScraper.Articles-contentScraper.Failures, contentScraper.Failures, time.Since(start), path)

		return nil
	},
}

func init() {
	cmdScrape.Args = cobra.ExactArgs(1)
	cmdScrape.PersistentFlags().StringVarP(&scrapeOutDir, "out", "o", "out/content/", "Directory where to store the scraped articles")
	cmdScrape.PersistentFlags().StringVarP(&timezone, "timezone", "t", "Europe/Berlin", "Timezone for storing the feeds")
	cmdScrape.PersistentFlags().BoolVarP(&scrapeVerbose, "verbose", "v", false, "Verbose logging of scraper")
	RootCmd.AddCommand(cmdScrape)
}

func getFeedsFilePath(itemsInputFile string, location *time.Location) (string, error) {
	if itemsInputFile == "" {
		return "", errors.New("Please provide a file with articles")
	}

	stat, err := os.Stat(itemsInputFile)
	if err != nil {
		return "", err
	}

	// Append current day to path when only received directory as input location
	if stat.IsDir() {
		day := time.Now().In(location)
		dayStr := day.Format("2-1-2006")
		path := filepath.Join(itemsInputFile, dayStr+".json")
		_, err := os.Stat(path)
		if err != nil {
			return path, nil
		}
	}

	return itemsInputFile, nil
}
