package crawler

import (
	"fmt"
	"time"

	"encoding/json"

	"io/ioutil"

	"github.com/mmcdole/gofeed"
)

// Feed represent an RSS/Atom feed
type Feed struct {
	URL   string      `json:"url"`
	Items []*FeedItem `json:"items"`
}

// FeedItem stores info of feed entry
type FeedItem struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	URL       string `json:"url"`
	Published string `json:"published"`
	GUID      string `json:"guid"`
}

// ScrapeFeeds download a list of provided feeds
func ScrapeFeeds(sources []string) error {
	feeds, err := fetch(sources)
	if err != nil {
		return err
	}

	err = save(feeds)
	return err
}

func fetch(sources []string) ([]Feed, error) {
	var feeds = make([]Feed, 0)
	for _, url := range sources {
		items, err := parse(url)
		if err != nil {
			return nil, err
		}

		feed := Feed{
			URL:   url,
			Items: items,
		}
		feeds = append(feeds, feed)
		fmt.Println(url)
	}
	return feeds, nil
}

func parse(url string) ([]*FeedItem, error) {
	feedParser := gofeed.NewParser()
	feed, err := feedParser.ParseURL(url)

	if err != nil {
		return nil, err
	}

	items := make([]*FeedItem, len(feed.Items))
	for i, item := range feed.Items {
		var GUID = item.GUID

		if GUID == "" {
			GUID = item.Link
		}
		newItem := FeedItem{
			Title:     item.Title,
			Content:   item.Description,
			URL:       item.Link,
			Published: item.Published,
			GUID:      item.GUID,
		}

		items[i] = &newItem
	}

	return items, nil
}

func save(feeds []Feed) error {
	// TODO update instead of overide file
	day := time.Now().Format("2-1-2006")
	filename := "out/" + day + ".json"
	jsonFeeds, err := json.Marshal(feeds)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, jsonFeeds, 0644)
}
