package cmd

import (
	"context"
	"fmt"
	"reflect"

	"github.com/spf13/cobra"
	"github.com/thesoenke/news-crawler/feedreader"
	"github.com/thesoenke/news-crawler/scraper"
	elastic "gopkg.in/olivere/elastic.v5"
)

var cmdNoD = &cobra.Command{
	Use: "nod",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := loadArticles()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(cmdNoD)
}

func loadArticles() ([]feedreader.FeedItem, error) {
	client, err := scraper.NewElasticClient()
	if err != nil {
		return nil, err
	}

	query := elastic.NewRangeQuery("published").From("now-1M").To("now")
	searchResult, err := client.Search().
		Index("news").
		Query(query).
		From(0).Size(10000).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	var feedItem feedreader.FeedItem
	for _, item := range searchResult.Each(reflect.TypeOf(feedItem)) {
		if article, ok := item.(feedreader.FeedItem); ok {
			fmt.Printf("Titel: %s\n", article.Title)
		}
	}

	return nil, nil
}
