package cmd

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/thesoenke/news-crawler/feedreader"
)

var cmdFeeds = &cobra.Command{
	Use:   "feeds",
	Short: "Scrape all provided feeds",
	RunE: func(cmd *cobra.Command, args []string) error {
		if inputFile == "" {
			return errors.New("Please provide a file with feeds")
		}

		reader, err := feedreader.New(inputFile)
		if err != nil {
			return err
		}

		start := time.Now()
		err = reader.Fetch()
		if err != nil {
			return err
		}

		items := 0
		for _, feed := range reader.Feeds {
			items += len(feed.Items)
		}
		log.Printf("Downloaded %d feeds with %d items in %s", len(reader.Feeds), items, time.Since(start))

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
	cmdFeeds.PersistentFlags().StringVarP(&inputFile, "file", "f", "feeds/feeds_de.txt", "Path to a file with feeds")
	cmdFeeds.PersistentFlags().StringVarP(&timezone, "timezone", "t", "Europe/Berlin", "Timezone for storing the feeds")
	cmdFeeds.PersistentFlags().StringVarP(&outDir, "out", "o", "out/feeds", "Directory where to store the feed items")
	RootCmd.AddCommand(cmdFeeds)
}
