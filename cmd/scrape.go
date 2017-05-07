package cmd

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/thesoenke/news-crawler/crawler"
)

var feedListFile string

var cmdScrape = &cobra.Command{
	Use:   "scrape",
	Short: "Scrape all provided articles",
	RunE: func(cmd *cobra.Command, args []string) error {
		if feedListFile == "" {
			return errors.New("Please provide a file with articles")
		}

		urls, err := extractArticleURLs(feedListFile)
		if err != nil {
			return err
		}
		return crawler.ScrapeURLs(urls)
	},
}

func init() {
	cmdScrape.PersistentFlags().StringVarP(&feedListFile, "articles", "a", "", "Path to a JSON file with feed items")
	cmdScrape.PersistentFlags().StringVarP(&timezone, "timezone", "t", "Europe/Berlin", "Timezone for storing the feeds")
	cmdScrape.PersistentFlags().StringVarP(&outDir, "out", "o", "out/", "Directory where to store the articles")
	RootCmd.AddCommand(cmdScrape)
}

func extractArticleURLs(path string) ([]string, error) {
	articlesFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var feeds []crawler.Feed
	err = json.Unmarshal(articlesFile, &feeds)
	if err != nil {
		return nil, err
	}

	var urls []string
	for _, feed := range feeds {
		for _, feedItem := range feed.Items {
			urls = append(urls, feedItem.URL)
		}
	}

	return urls, nil
}
