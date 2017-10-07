package cmd

import (
	"fmt"
	"path"
	"time"

	"github.com/spf13/cobra"
	"github.com/thesoenke/news-crawler/feedreader"
)

var feedsOutDir string
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
		err = reader.Fetch(&dayTime, verbose)
		if err != nil {
			return err
		}

		items := 0
		for _, feed := range reader.Feeds {
			items += len(feed.Items)
		}

		err = reader.LogFailures(feedsOutDir, &dayTime)
		if err != nil {
			return err
		}

		dir := path.Join(feedsOutDir, lang)
		file, err := reader.Store(dir, &dayTime)
		if err != nil {
			return err
		}

		successLog := fmt.Sprintf("Feedreader %s\nSuccessful: %d\nFailures: %d\nArticles: %d\nTime: %s\nFile: %s\n\n", time.Now().In(location), len(reader.Feeds), len(reader.FailedFeeds), items, time.Since(start), file)
		fmt.Println(successLog)
		err = writeLog(logsDir, successLog)
		return err
	},
}

func init() {
	cmdFeeds.Args = cobra.ExactArgs(1)
	cmdFeeds.PersistentFlags().StringVarP(&lang, "lang", "l", "", "Language of the content")
	cmdFeeds.PersistentFlags().StringVarP(&timezone, "timezone", "t", "Europe/Berlin", "Timezone for storing the feeds")
	cmdFeeds.PersistentFlags().StringVarP(&feedsOutDir, "dir", "d", "out/feeds", "Directory to store feed items")
	cmdFeeds.PersistentFlags().StringVar(&logsDir, "logs", "out/log", "File to store logs")
	cmdFeeds.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Output more detailed logging")
	RootCmd.AddCommand(cmdFeeds)
}
