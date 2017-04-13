package main

import (
	"github.com/mmcdole/gofeed"
)

// Article is a single feed item
type Article struct {
	Title     string
	Content   string
	URL       string
	Published string
}

// Parse a feed passed as a the URL
func Parse(url string) ([]Article, error) {
	fp := gofeed.NewParser()

	feed, err := fp.ParseURL("http://feeds.twit.tv/twit.xml")

	if err != nil {
		return nil, err
	}

	articles := make([]Article, 0)
	for _, item := range feed.Items {
		article := Article{
			Title:     item.Title,
			Content:   item.Description,
			URL:       item.Link,
			Published: item.Published,
		}

		articles = append(articles, article)
	}

	return articles, nil
}
