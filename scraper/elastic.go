package scraper

import (
	"context"
	"os"

	elastic "gopkg.in/olivere/elastic.v5"
)

func NewElasticClient() (*elastic.Client, error) {
	elasticURL := elastic.SetURL(os.Getenv("ELASTIC_URL"))
	auth := elastic.SetBasicAuth(os.Getenv("ELASTIC_USER"), os.Getenv("ELASTIC_PASSWORD"))
	client, err := elastic.NewClient(elasticURL, auth)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (article *Article) StoreElastic(client *elastic.Client) error {
	ctx := context.Background()
	_, err := client.Index().
		Index("news").
		Type("article").
		BodyJson(article.FeedItem).
		Refresh("true").
		Do(ctx)

	if err != nil {
		return err
	}

	return nil
}

func createIndex(client *elastic.Client) error {
	ctx := context.Background()
	exists, err := client.IndexExists("news").Do(ctx)
	if err != nil {
		return err
	}

	if !exists {
		_, err := client.CreateIndex("news").Do(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
