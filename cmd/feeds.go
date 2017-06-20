package cmd

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/thesoenke/news-crawler/feedreader"
)

var feedsFile string
var timezone string
var outDir string

var cmdFeeds = &cobra.Command{
	Use:   "feeds",
	Short: "Scrape all provided feeds",
	RunE: func(cmd *cobra.Command, args []string) error {
		if feedsFile == "" {
			return errors.New("Please provide a file with feeds")
		}

		reader, err := feedreader.New(feedsFile)
		if err != nil {
			return err
		}

		start := time.Now()
		_, err = reader.Fetch()
		if err != nil {
			return err
		}
		log.Printf("Feedreader finished in %s", time.Since(start))

		location, err := time.LoadLocation(timezone)
		if err != nil {
			return err
		}

		err = reader.Store(outDir, location)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	cmdFeeds.PersistentFlags().StringVarP(&feedsFile, "feeds", "f", "feeds/feeds_de.txt", "Path to a file with feeds")
	cmdFeeds.PersistentFlags().StringVarP(&timezone, "timezone", "t", "Europe/Berlin", "Timezone for storing the feeds")
	cmdFeeds.PersistentFlags().StringVarP(&outDir, "out", "o", "out/", "Directory where to store the feed items")
	RootCmd.AddCommand(cmdFeeds)
}
