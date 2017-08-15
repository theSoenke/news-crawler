package feedreader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// Store all downloaded feeds into a JSON file
func (fr *FeedReader) Store(outDir string, dayTime *time.Time) error {
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		err := os.MkdirAll(outDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	day := dayTime.Format("2-1-2006")
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

// LogFailures stores all failed feed downloads
func (fr *FeedReader) LogFailures(dir string, dayTime *time.Time) error {
	logDir := filepath.Join(dir, "log")
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	filename := filepath.Join(logDir, "failed.txt")
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	logText := ""
	for _, url := range fr.FailedFeeds {
		logText += fmt.Sprintf("%s,%s", dayTime, url+"\n")
	}

	_, err = file.WriteString(logText)
	if err != nil {
		return err
	}
	err = file.Close()
	return err
}
