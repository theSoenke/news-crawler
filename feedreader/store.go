package feedreader

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func (fr *FeedReader) Store(outDir string, location *time.Location) error {
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		err := os.MkdirAll(outDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	dayLocation := time.Now().In(location)
	day := dayLocation.Format("2-1-2006")
	feedFile := filepath.Join(outDir, day+".json")
	feeds := fr.Feeds
	if _, err := os.Stat(feedFile); !os.IsNotExist(err) {
		feedsFile, err := ioutil.ReadFile(feedFile)
		if err != nil {
			return err
		}

		var oldFeeds []Feed
		err = json.Unmarshal(feedsFile, &oldFeeds)
		if err != nil {
			return err
		}

		feeds = merge(feeds, oldFeeds)
	}

	jsonFeeds, err := json.Marshal(feeds)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(feedFile, jsonFeeds, 0644)
	return err
}
