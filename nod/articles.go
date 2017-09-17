package nod

import (
	"bytes"
	"context"
	"fmt"
	"reflect"

	"github.com/thesoenke/news-crawler/feedreader"
	"github.com/thesoenke/news-crawler/scraper"
	elastic "gopkg.in/olivere/elastic.v5"
)

func CreateNoDInput() error {
	output, err := outputText()
	if err != nil {
		return err
	}

	fmt.Println(output)
	return nil
}

func compressOutput(output string) {
	// TODO
}

func outputText() (string, error) {
	from := "2017-08-01"
	to := "2017-09-30"
	articles, err := loadArticles(from, to)
	if err != nil {
		return "", nil
	}

	tokennizer, err := NewSentenceTokenizer("german")
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	for _, article := range articles {
		sentences := tokennizer.Tokenize(article.Content)
		for _, s := range sentences {
			output := fmt.Sprintf("%s\t%s\n", s.Text, article.URL)
			buffer.WriteString(output)
		}
	}

	return buffer.String(), nil
}

func loadArticles(from string, to string) ([]feedreader.FeedItem, error) {
	client, err := scraper.NewElasticClient()
	if err != nil {
		return nil, err
	}

	query := elastic.NewRangeQuery("published").From(from).To(to)
	searchResult, err := client.Search().
		Index("news").
		Query(query).
		From(0).Size(10000).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	items := make([]feedreader.FeedItem, 0)
	var feedItem feedreader.FeedItem
	for _, item := range searchResult.Each(reflect.TypeOf(feedItem)) {
		if article, ok := item.(feedreader.FeedItem); ok {
			// fmt.Printf("Titel: %s\n", article.Title)
			items = append(items, article)
		}
	}

	return items, nil
}
