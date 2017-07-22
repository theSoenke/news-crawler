package feedreader

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/SlyMarbo/rss"
	"gopkg.in/cheggaaa/pb.v1"
)

const (
	userAgent = "Mozilla/5.0 (X11; Linux x86_64; rv:53.0) Gecko/20100101 Firefox/53.0"
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
	Sources     []string
	Feeds       []Feed
	FailedFeeds []string
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

// Fetch feed items
func (fr *FeedReader) Fetch(verbose bool) error {
	concurrencyLimit := 100
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
	failedFeeds := make([]string, 0)
	for i := 0; i < count; i++ {
		select {
		case feed := <-feedChan:
			feeds = append(feeds, *feed)
		case url := <-errURLChan:
			failedFeeds = append(failedFeeds, url)
		}
		bar.Increment()
	}
	bar.Finish()

	fr.Feeds = feeds
	fr.FailedFeeds = failedFeeds

	return nil
}

func fetchFeedURL(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	timeout := time.Duration(60 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	return client.Do(req)
}

func fetchFeed(url string) ([]*FeedItem, error) {
	feed, err := rss.FetchByFunc(fetchFeedURL, url)
	if err != nil {
		return nil, err
	}

	items := make([]*FeedItem, 0)
	for _, item := range feed.Items {
		newItem := FeedItem{
			Title:     item.Title,
			Content:   item.Content,
			URL:       item.Link,
			Published: item.Date.String(),
			GUID:      item.ID,
		}

		err := newItem.validate()
		if err != nil {
			continue
		}

		items = append(items, &newItem)
	}

	return items, nil
}

func (item *FeedItem) validate() error {
	if item.URL == "" {
		return fmt.Errorf("Feed item contains no url: %s", item)
	}

	if item.GUID == "" {
		item.GUID = item.URL
	}

	return nil
}
