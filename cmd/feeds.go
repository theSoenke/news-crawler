package cmd

import (
	"bufio"
	"errors"
	"os"
	"time"

	"fmt"

	"github.com/spf13/cobra"
	"github.com/thesoenke/news-crawler/crawler"
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

		location, err := time.LoadLocation(timezone)
		if err != nil {
			return err
		}

		feeds, err := loadFeeds(feedsFile)
		if err != nil {
			return err
		}

		fmt.Print(feeds)

		err = crawler.ScrapeFeeds(feeds, outDir, location)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	cmdFeeds.PersistentFlags().StringVarP(&feedsFile, "feeds", "f", "feeds/news_de.json", "Path to a JSON file with feeds")
	cmdFeeds.PersistentFlags().StringVarP(&timezone, "timezone", "t", "Europe/Berlin", "Timezone for storing the feeds")
	cmdFeeds.PersistentFlags().StringVarP(&outDir, "out", "o", "out/", "Directory where to store the feed items")
	RootCmd.AddCommand(cmdFeeds)
}

func loadFeeds(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	feeds := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		feeds = append(feeds, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return feeds, nil
}
