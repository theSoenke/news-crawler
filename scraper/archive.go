package scraper

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
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

	hash := md5.Sum([]byte(article.FeedItem.URL))
	filename := hex.EncodeToString(hash[:]) + ".html"
	articlePath := filepath.Join(dayDir, filename)
	err := ioutil.WriteFile(articlePath, []byte(article.HTML), 0644)
	return err
}
