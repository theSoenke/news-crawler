package scraper

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Write article to file
func (article *Article) Write(outDir string, dayTime *time.Time) error {
	day := dayTime.Format("2-1-2006")
	dayDir := filepath.Join(outDir, day)

	if _, err := os.Stat(dayDir); os.IsNotExist(err) {
		err := os.MkdirAll(dayDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	filename := strings.Replace(filepath.FromSlash(article.FeedItem.URL), "/", "\\", -1) + ".html"
	articlePath := filepath.Join(dayDir, filename)
	err := ioutil.WriteFile(articlePath, []byte(article.HTML), 0644)
	if err != nil {
		return err
	}

	return nil
}
