package cmd

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/thesoenke/news-crawler/scraper"
)

var feedListFile string

var cmdScrape = &cobra.Command{
	Use:   "scrape",
	Short: "Scrape all provided articles",
	RunE: func(cmd *cobra.Command, args []string) error {
		if feedListFile == "" {
			return errors.New("Please provide a file with articles")
		}

		contentScraper, err := scraper.New(feedListFile)
		if err != nil {
			return err
		}

		start := time.Now()
		err = contentScraper.Scrape()
		if err != nil {
			return err
		}

		articles := 0
		for _, feed := range contentScraper.Feeds {
			articles += len(feed.Items)
		}
		log.Printf("Scraper downloaded %d articles in %s", articles, time.Since(start))

		return nil
	},
}

func init() {
	cmdScrape.PersistentFlags().StringVarP(&feedListFile, "articles", "a", "", "Path to a JSON file with feed items")
	cmdScrape.PersistentFlags().StringVarP(&timezone, "timezone", "t", "Europe/Berlin", "Timezone for storing the feeds")
	cmdScrape.PersistentFlags().StringVarP(&outDir, "out", "o", "out/", "Directory where to store the articles")
	RootCmd.AddCommand(cmdScrape)
}
