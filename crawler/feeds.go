package crawler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

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
func ScrapeFeeds(sources []string, outDir string, timezone *time.Location) error {
	feeds, err := fetch(sources)
	if err != nil {
		return err
	}

	err = save(feeds, outDir, timezone)
	return err
}

func fetch(sources []string) ([]Feed, error) {
	var feeds = make([]Feed, 0)
	for i, url := range sources {
		items, err := parse(url)
		if err != nil {
			fmt.Printf("%d: Failed %s\n", i, url)
			continue
		}

		feed := Feed{
			URL:   url,
			Items: items,
		}
		feeds = append(feeds, feed)
		fmt.Printf("%d: %s\n", i, url)
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
		newItem := FeedItem{
			Title:     item.Title,
			Content:   item.Description,
			URL:       item.Link,
			Published: item.Published,
			GUID:      item.GUID,
		}

		err := newItem.validate()
		if err != nil {
			return nil, err
		}

		items[i] = &newItem
	}
	return items, nil
}

func save(feeds []Feed, outDir string, location *time.Location) error {
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		err := os.MkdirAll(outDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	dayLocation := time.Now().In(location)
	day := dayLocation.Format("2-1-2006")
	feedFile := outDir + day + ".json"
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
	return ioutil.WriteFile(feedFile, jsonFeeds, 0644)
}

func merge(newFeeds []Feed, oldFeeds []Feed) []Feed {
	oldFeedsMap := make(map[string][]*FeedItem, len(oldFeeds))
	feeds := make([]Feed, len(newFeeds))

	for _, feed := range oldFeeds {
		oldFeedsMap[feed.URL] = feed.Items
	}

	for _, feed := range newFeeds {
		if oldFeedsMap[feed.URL] != nil {
			items := removeDuplicates(feed.Items, oldFeedsMap[feed.URL])
			newFeed := Feed{
				URL:   feed.URL,
				Items: items,
			}
			feeds = append(feeds, newFeed)
		} else {
			feeds = append(feeds, feed)
		}
	}

	return feeds
}

func removeDuplicates(newItems []*FeedItem, oldItems []*FeedItem) []*FeedItem {
	found := make(map[string]bool)

	for _, item := range oldItems {
		found[item.GUID] = true
	}

	for _, item := range newItems {
		if !found[item.GUID] {
			oldItems = append(oldItems, item)
		}
	}

	return oldItems
}

func (item *FeedItem) validate() error {
	if item.URL == "" {
		return fmt.Errorf("Feed item contains no url: %s", item)
	}

	if item.GUID == "" {
		item.GUID = item.URL
	}

	if item.Published == "" {
		return fmt.Errorf("Feed item contains no published date: %s", item.URL)
	}

	return nil
}
