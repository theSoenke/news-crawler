package scraper

import (
	"bytes"
	"fmt"
	"net/url"
	"unicode/utf8"

	"github.com/advancedlogic/GoOse"
	"github.com/jlubawy/go-boilerpipe"
	"github.com/jlubawy/go-boilerpipe/extractor"
)

// Extract the content of an article
func (article *Article) Extract() error {
	content, err := ExtractContentGoOseUtf8(article.FeedItem.URL, article.HTML)
	// content, err := ExtractContentBoilerpipe(article.FeedItem.URL, article.HTML)
	if err != nil {
		return err
	}

	article.FeedItem.Content = content
	return nil
}

var issues int

func ExtractContentGoOse(url string, html string) (string, error) {
	validUTF8 := utf8.Valid([]byte(html))
	if !validUTF8 {
		issues++
		fmt.Println(issues)
		// fmt.Println("No UTF8")
	}
	g := goose.New()
	article, err := g.ExtractFromRawHTML(url, html)
	if err != nil {
		return "", err
	}

	return article.CleanedText, nil
}

func ExtractContentGoOseUtf8(url string, html string) (string, error) {
	g := goose.New()
	article, err := g.ExtractFromRawHTML(url, html)
	if err != nil {
		return "", err
	}

	// goose.UTF8encode(article.CleanedText, sourceCharset string)

	validUTF8 := utf8.Valid([]byte(article.CleanedText))
	if !validUTF8 {
		fmt.Println("No UTF8")
	}

	return article.CleanedText, nil
}

func ExtractContentBoilerpipe(urlStr string, html string) (string, error) {
	validUTF8 := utf8.Valid([]byte(html))
	if !validUTF8 {
		fmt.Println("No UTF8")
	}

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
