package cmd

import (
	"encoding/json"
	"io/ioutil"

	"log"

	"github.com/spf13/cobra"
	"github.com/thesoenke/news-crawler/crawler"
)

var feedListFile string

var cmdScrape = &cobra.Command{
	Use:   "scrape",
	Short: "Scrape all provided articles",
	Run: func(cmd *cobra.Command, args []string) {
		urls, err := extractArticleURLs(feedListFile)
		if err != nil {
			log.Fatal(err)
		}
		crawler.ScrapeURLs(urls)
	},
}

func init() {
	cmdScrape.PersistentFlags().StringVarP(&feedListFile, "articles", "a", "feeds/news_de.json", "Path to a JSON file with feeds")
	cmdScrape.PersistentFlags().StringVarP(&timezone, "timezone", "t", "Europe/Berlin", "Timezone for storing the feeds")
	cmdScrape.PersistentFlags().StringVarP(&outDir, "out", "o", "out/", "Directory where to store the feed items")
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
