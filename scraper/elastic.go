package scraper

import (
	"context"
	"os"

	elastic "gopkg.in/olivere/elastic.v5"
)

// NewElasticClient creates a new client to connect to an elasticsearch cluster
func NewElasticClient() (*elastic.Client, error) {
	elasticURL := elastic.SetURL(os.Getenv("ELASTIC_URL"))
	auth := elastic.SetBasicAuth(os.Getenv("ELASTIC_USER"), os.Getenv("ELASTIC_PASSWORD"))
	client, err := elastic.NewClient(elasticURL, auth)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Index article in elasticsearch
func (scraper *Scraper) index(article *Article) error {
	ctx := context.Background()
	_, err := scraper.ElasticClient.Index().
		Index("news").
		Type("article").
		BodyJson(article.FeedItem).
		Refresh("true").
		Do(ctx)
	return err
}

func (scraper *Scraper) createIndex() error {
	ctx := context.Background()
	exists, err := scraper.ElasticClient.IndexExists("news").Do(ctx)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	_, err = scraper.ElasticClient.CreateIndex("news").Do(ctx)
	return err
}
