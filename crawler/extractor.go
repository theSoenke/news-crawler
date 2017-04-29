package crawler

import (
	"github.com/advancedlogic/GoOse"
)

func extract(url string, html string) (string, error) {
	g := goose.New()
	article, err := g.ExtractFromRawHTML(url, html)
	if err != nil {
		return "", err
	}

	return article.CleanedText, nil
}
