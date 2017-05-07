package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/thesoenke/news-crawler/crawler"
)

var feeds string
var timezone string
var outDir string

var cmdFeeds = &cobra.Command{
	Use:   "feeds",
	Short: "Scrape all provided feeds",
	Run: func(cmd *cobra.Command, args []string) {
		location, err := time.LoadLocation(timezone)
		if err != nil {
			log.Fatal("Please provide a valid timezone with --timezone")
		}

		sources, err := extractFeedURLs(feeds)
		if err != nil {
			log.Fatal("Failed to import sources")
		}

		err = crawler.ScrapeFeeds(sources, outDir, location)
		if err != nil {
			log.Fatal(err.Error())
		}
	},
}

func init() {
	cmdFeeds.PersistentFlags().StringVarP(&feeds, "feeds", "f", "feeds/news_de.json", "Path to a JSON file with feeds")
	cmdFeeds.PersistentFlags().StringVarP(&timezone, "timezone", "t", "Europe/Berlin", "Timezone for storing the feeds")
	cmdFeeds.PersistentFlags().StringVarP(&outDir, "out", "o", "out/", "Directory where to store the feed items")
	RootCmd.AddCommand(cmdFeeds)
}

func extractFeedURLs(path string) ([]string, error) {
	sourceFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var sources = make([]string, 0)
	err = json.Unmarshal(sourceFile, &sources)
	if err != nil {
		return nil, err
	}

	return sources, nil
}
