package crawler

import (
	"fmt"

	"github.com/mmcdole/gofeed"
)

// Feed represent an RSS/Atom feed
type Feed struct {
	URL   string
	Items []*FeedItem
}

// FeedItem stores info of feed entry
type FeedItem struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	URL       string `json:"url"`
	Published string `json:"published"`
	GUID      string `json:"guid"`
}

// Run a new crawler
func Run(sources []string) ([]Feed, error) {
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
			return nil, fmt.Errorf("GUID is empty. Link: %s", item.Link)
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
