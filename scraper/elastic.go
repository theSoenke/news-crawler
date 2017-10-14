package scraper

import (
	"context"
	"fmt"
	"os"

	elastic "gopkg.in/olivere/elastic.v5"
)

const (
	ElasticHost     = "http://localhost"
	ElasticPort     = 9200
	ElasticUser     = "elastic"
	ElasticPassword = "changeme"
)

// NewElasticClient creates a new client to connect to an elasticsearch cluster
func NewElasticClient() (*elastic.Client, error) {
	url := os.Getenv("ELASTIC_URL")
	user := os.Getenv("ELASTIC_USER")
	password := os.Getenv("ELASTIC_PASSWORD")

	if url == "" {
		url = fmt.Sprintf("%s:%d", ElasticHost, ElasticPort)
	}
	if user == "" {
		user = ElasticUser
	}
	if password == "" {
		user = ElasticPassword
	}

	auth := elastic.SetBasicAuth(user, password)
	client, err := elastic.NewClient(elastic.SetURL(url), auth)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Index article in elasticsearch
func (scraper *Scraper) index(article *Article) error {
	ctx := context.Background()
	_, err := scraper.ElasticClient.Index().
		Index("news-" + scraper.Lang).
		Type("article").
		BodyJson(article.FeedItem).
		Refresh("true").
		Do(ctx)
	return err
}

// logError in ElasticSearch
func (scraper *Scraper) logError(fetchError *FetchError) error {
	ctx := context.Background()
	_, err := scraper.ElasticClient.Index().
		Index("failures-" + scraper.Lang).
		Type("failure").
		BodyJson(fetchError).
		Refresh("true").
		Do(ctx)
	return err
}
