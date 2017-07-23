package cmd

import (
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/thesoenke/news-crawler/feedreader"
)

var feedOutDir string
var feedsVerbose bool

var cmdFeeds = &cobra.Command{
	Use:   "feeds",
	Short: "Scrape all provided feeds",
	RunE: func(cmd *cobra.Command, args []string) error {
		feedInputFile := args[0]
		reader, err := feedreader.New(feedInputFile)
		if err != nil {
			return err
		}

		location, err := time.LoadLocation(timezone)
		if err != nil {
			return err
		}

		dayTime := time.Now().In(location)
		start := time.Now()
		err = reader.Fetch(&dayTime, feedsVerbose)
		if err != nil {
			return err
		}

		items := 0
		for _, feed := range reader.Feeds {
			items += len(feed.Items)
		}
		log.Printf("Feeds: %d successful, %d failures, %d items in %s", len(reader.Feeds), len(reader.FailedFeeds), items, time.Since(start))

		err = reader.LogFailures(feedOutDir, &dayTime)
		if err != nil {
			return err
		}

		err = reader.Store(feedOutDir, &dayTime)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	cmdFeeds.Args = cobra.ExactArgs(1)
	cmdFeeds.PersistentFlags().StringVarP(&timezone, "timezone", "t", "Europe/Berlin", "Timezone for storing the feeds")
	cmdFeeds.PersistentFlags().StringVarP(&feedOutDir, "out", "o", "out/feeds/", "Directory where to store the feed items")
	cmdFeeds.PersistentFlags().BoolVarP(&feedsVerbose, "verbose", "v", false, "Output more detailed logging")
	RootCmd.AddCommand(cmdFeeds)
}
