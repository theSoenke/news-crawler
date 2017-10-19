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
func CreateCorpus(lang string, from string, timeZone string, dir string) error {
	day, err := time.Parse("2-1-2006", from)
	if err != nil {
		return err
	}

	tokenizer, err := NewSentenceTokenizer(lang)
	if err != nil {
		return err
	}

	for {
		today := time.Now().In(day.Location()).Truncate(time.Hour * 24)
		if day.After(today) || day.Equal(today) {
			break
		}

		corpus := dayCorpus{Day: day}
		output, err := corpus.generate(lang, timeZone, tokenizer)
		if err != nil {
			return err
		}

		if output == "" {
			day = day.AddDate(0, 0, 1)
			continue
		}

		err = corpus.compress(output, dir, day.Format("20060102"))
		if err != nil {
			return err
		}

		fmt.Println(day.Format("2006-01-02"))
		day = day.AddDate(0, 0, 1)
	}

	return nil
}

func (corpus *dayCorpus) generate(lang string, timeZone string, tokenizer sentences.SentenceTokenizer) (string, error) {
	day := corpus.Day.Format("2006-01-02")
	articles, err := loadArticles(lang, day, day, timeZone)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	for _, article := range articles {
		sentences := tokenizer.Tokenize(article.Content)
		for _, s := range sentences {
			text := strings.Join(strings.Fields(s.Text), " ")
			if len(text) < 20 || len(text) > 256 {
				continue
			}
			text = strings.Replace(text, "|", " ", -1)
			output := fmt.Sprintf("%s\t%s\n", text, article.URL)
			buffer.WriteString(output)
		}
	}

	return buffer.String(), nil
}

func loadArticles(lang string, from string, to string, timeZone string) ([]feedreader.FeedItem, error) {
	client, err := scraper.NewElasticClient()
	if err != nil {
		return nil, err
	}

	query := elastic.NewRangeQuery("published").
		TimeZone(timeZone).
		From(from).
		To(to)

	searchResult, err := client.Search().
		Index("news-" + lang).
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
