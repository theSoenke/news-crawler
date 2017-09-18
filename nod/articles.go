package nod

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/thesoenke/news-crawler/feedreader"
	"github.com/thesoenke/news-crawler/scraper"
	elastic "gopkg.in/olivere/elastic.v5"
)

func CreateNoDCorpus() error {
	startDay := "02-08-2017" // TODO get day of first article from Elasticsearch
	day, err := time.Parse("2-1-2006", startDay)
	if err != nil {
		return err
	}

	for {
		day = day.AddDate(0, 0, 1)
		if time.Now().AddDate(0, 0, -1).Before(day) {
			break
		}

		output, err := generateDayOutput(day)
		if err != nil {
			return err
		}

		err = compressBz2(output, day.Format("20060102"))
		if err != nil {
			return err
		}
	}

	return nil
}

func generateDayOutput(day time.Time) (string, error) {
	from := day.Format("2006-01-02")
	to := day.AddDate(0, 0, 0).Format("2006-01-02")
	articles, err := loadArticles(from, to)
	if err != nil {
		return "", err
	}

	tokennizer, err := NewSentenceTokenizer("german")
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	for _, article := range articles {
		sentences := tokennizer.Tokenize(article.Content)
		for _, s := range sentences {
			text := strings.Join(strings.Fields(s.Text), " ")
			if len(text) < 20 {
				continue
			}
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
			items = append(items, article)
		}
	}

	return items, nil
}
