package stats

import (
	"context"
	"fmt"

	"github.com/thesoenke/news-crawler/scraper"
)

func TotalArticles() error {
	client, err := scraper.NewElasticClient()
	if err != nil {
		return err
	}

	searchResult, err := client.Search().
		Index("news").
		From(0).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("Total articles in index %d\n", searchResult.Hits.TotalHits)
	return nil
}
