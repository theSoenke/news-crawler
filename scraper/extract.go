package scraper

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/advancedlogic/GoOse"
	"github.com/jlubawy/go-boilerpipe"
	"github.com/jlubawy/go-boilerpipe/extractor"
)

// Fetch the content of an article from the web
func (article *Article) Fetch() error {
	timeout := time.Duration(60 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest("GET", article.FeedItem.URL, nil)
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return fmt.Errorf("Site returned status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	article.HTML = string(body)

	return nil
}

// Extract the content of an article
func (article *Article) Extract() error {
	// content, err := extractContentGoOse(article.FeedItem.URL, article.HTML)
	content, err := extractContentBoilerpipe(article.FeedItem.URL, article.HTML)
	if err != nil {
		return err
	}

	article.FeedItem.Content = content

	return nil
}

func extractContentGoOse(url string, html string) (string, error) {
	g := goose.New()
	article, err := g.ExtractFromRawHTML(url, html)
	if err != nil {
		return "", err
	}

	return article.CleanedText, nil
}

func extractContentBoilerpipe(urlStr string, html string) (string, error) {
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
