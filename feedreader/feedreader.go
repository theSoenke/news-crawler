package feedreader

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/mmcdole/gofeed"
	"gopkg.in/cheggaaa/pb.v1"
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

type FeedReader struct {
	Sources    []string
	Feeds      []Feed
	FailedURLs []string
}

// New creates a feedreader
func New(feedsFile string) (FeedReader, error) {
	feedreader := FeedReader{}
	feeds, err := loadFeeds(feedsFile)
	if err != nil {
		return feedreader, err
	}

	feedreader.Sources = feeds

	return feedreader, nil
}

func (fr *FeedReader) Fetch(verbose bool) error {
	concurrencyLimit := 200
	wg := sync.WaitGroup{}
	queue := make(chan string)
	errURLChan := make(chan string)
	feedChan := make(chan *Feed)
	count := len(fr.Sources)
	bar := pb.StartNew(count)

	for worker := 0; worker < concurrencyLimit; worker++ {
		wg.Add(1)

		go func(verbose bool) {
			defer wg.Done()

			for url := range queue {
				items, err := fetchFeed(url)
				if err != nil {
					if verbose {
						log.Printf("Failed to fetch feed %s %s", url, err)
					}
					errURLChan <- url
					continue
				}

				feed := Feed{
					URL:   url,
					Items: items,
				}
				feedChan <- &feed
			}
		}(verbose)
	}

	go func(feeds []string) {
		for _, url := range feeds {
			queue <- url
		}
	}(fr.Sources)

	feeds := make([]Feed, 0)
	failed := make([]string, 0)
	for i := 0; i < count; i++ {
		select {
		case feed := <-feedChan:
			feeds = append(feeds, *feed)
		case url := <-errURLChan:
			failed = append(failed, url)
		}
		bar.Increment()
	}
	bar.Finish()

	fr.Feeds = feeds
	fr.FailedURLs = failed

	return nil
}

func fetchFeed(url string) ([]*FeedItem, error) {
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

func loadFeeds(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	feeds := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		feeds = append(feeds, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return feeds, nil
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
		item.Published = time.Now().String()
	}

	return nil
}
