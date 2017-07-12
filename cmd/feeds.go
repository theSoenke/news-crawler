package cmd

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/thesoenke/news-crawler/feedreader"
)

var feedInputFile string
var feedOutDir string
var feedsVerbose bool

var cmdFeeds = &cobra.Command{
	Use:   "feeds",
	Short: "Scrape all provided feeds",
	RunE: func(cmd *cobra.Command, args []string) error {
		if feedInputFile == "" {
			return errors.New("Please provide a file with feeds")
		}

		reader, err := feedreader.New(feedInputFile)
		if err != nil {
			return err
		}

		start := time.Now()
		err = reader.Fetch(feedsVerbose)
		if err != nil {
			return err
		}

		items := 0
		for _, feed := range reader.Feeds {
			items += len(feed.Items)
		}
		log.Printf("Feeds: %d successful, %d failures, %d items in %s", len(reader.Feeds)-len(reader.FailedURLs), len(reader.FailedURLs), items, time.Since(start))

		location, err := time.LoadLocation(timezone)
		if err != nil {
			return err
		}

		err = reader.LogFailures(feedOutDir, location)
		if err != nil {
			return err
		}

		err = reader.Store(feedOutDir, location)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	cmdFeeds.PersistentFlags().StringVarP(&feedInputFile, "file", "f", "feeds/feeds_de.txt", "Path to a file with feeds")
	cmdFeeds.PersistentFlags().StringVarP(&timezone, "timezone", "t", "Europe/Berlin", "Timezone for storing the feeds")
	cmdFeeds.PersistentFlags().StringVarP(&feedOutDir, "out", "o", "out/feeds/", "Directory where to store the feed items")
	cmdFeeds.PersistentFlags().BoolVarP(&feedsVerbose, "verbose", "v", false, "Output more detailed logging")
	RootCmd.AddCommand(cmdFeeds)
}
