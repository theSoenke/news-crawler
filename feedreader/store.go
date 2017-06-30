package feedreader

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type FeedFailure struct {
	URL   string `json:"url"`
	Count int    `json:"count"`
}

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

func (fr *FeedReader) StoreFailures(dir string) error {
	path := filepath.Join(dir, "failures.json")
	failuresJSON, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var oldFailures []FeedFailure
	err = json.Unmarshal(failuresJSON, &oldFailures)

	newFailures := make([]FeedFailure, len(fr.FailedURLs))
	for _, url := range fr.FailedURLs {
		failure := FeedFailure{
			URL:   url,
			Count: 1,
		}
		newFailures = append(newFailures, failure)
	}

	failures := mergeFailures(newFailures, oldFailures)
	failuresJSON, err = json.Marshal(failures)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, failuresJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}

func mergeFailures(newFailures []FeedFailure, oldFailures []FeedFailure) []FeedFailure {
	failureURLs := make(map[string]int, len(oldFailures))
	for _, failure := range oldFailures {
		failureURLs[failure.URL] = failure.Count
	}

	for _, failure := range newFailures {
		failureURLs[failure.URL] += failure.Count
	}

	failures := make([]FeedFailure, 0)
	for url, count := range failureURLs {
		failure := FeedFailure{
			URL:   url,
			Count: count,
		}
		failures = append(failures, failure)
	}

	return failures
}
