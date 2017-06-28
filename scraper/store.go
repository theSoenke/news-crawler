package scraper

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

func (scraper *Scraper) Store(outDir string, location *time.Location) error {
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		err := os.MkdirAll(outDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// TODO support storage of articles
	feedsJSON, err := json.Marshal(scraper.Feeds)
	if err != nil {
		return err
	}

	dayLocation := time.Now().In(location)
	day := dayLocation.Format("2-1-2006")
	contentFile := outDir + day + ".json"
	err = ioutil.WriteFile(contentFile, feedsJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}
