package scraper

import (
	"bytes"
	"net/url"

	"github.com/jlubawy/go-boilerpipe"
	"github.com/jlubawy/go-boilerpipe/extractor"
	goose "github.com/thesoenke/GoOse"
)

// Extract the content of an article
func (article *Article) Extract() error {
	content, err := ExtractWithGoOse(article.FeedItem.URL, article.HTML)
	// content, err := ExtractWithBoilerpipe(article.FeedItem.URL, article.HTML)
	if err != nil {
		return err
	}

	article.FeedItem.Content = content
	return nil
}

func ExtractWithGoOse(url string, html string) (string, error) {
	g := goose.New()
	article, err := g.ExtractFromRawHTML(url, html)
	if err != nil {
		return "", err
	}

	return article.CleanedText, nil
}

func ExtractWithBoilerpipe(urlStr string, html string) (string, error) {
	reader := bytes.NewReader([]byte(html))
	url, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	doc, err := boilerpipe.NewDocument(reader, url)
	if err != nil {
		return "", err
	}
	extractor.Article().Process(doc)
	return doc.Content(), nil
}
