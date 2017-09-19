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
	"gopkg.in/neurosnap/sentences.v1"
	elastic "gopkg.in/olivere/elastic.v5"
)

type dayCorpus struct {
	Day time.Time
}

// CreateCorpus creates input for NoDCore from an ElasticSearch instance
func CreateCorpus(language string, dir string) error {
	startDay := "02-08-2017" // TODO get day of first article from ElasticSearch
	day, err := time.Parse("2-1-2006", startDay)
	if err != nil {
		return err
	}

	tokenizer, err := NewSentenceTokenizer(language)
	if err != nil {
		return err
	}

	for {
		day = day.AddDate(0, 0, 1)
		if time.Now().AddDate(0, 0, -1).Before(day) {
			break
		}

		corpus := dayCorpus{
			Day: day,
		}
		output, err := corpus.generate(tokenizer)
		if err != nil {
			return err
		}

		err = corpus.compress(output, dir, day.Format("20060102"))
		if err != nil {
			return err
		}
	}

	return nil
}

func (corpus *dayCorpus) generate(tokenizer sentences.SentenceTokenizer) (string, error) {
	from := corpus.Day.Format("2006-01-02")
	to := corpus.Day.AddDate(0, 0, 0).Format("2006-01-02")
	articles, err := loadArticles(from, to)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	for _, article := range articles {
		sentences := tokenizer.Tokenize(article.Content)
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
