package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"log"

	"github.com/thesoenke/news-crawler/feedreader"
	"gopkg.in/cheggaaa/pb.v1"
)

const (
	userAgent = "Mozilla/5.0 (X11; Linux x86_64; rv:53.0) Gecko/20100101 Firefox/53.0"
)

type Scraper struct {
	Feeds []feedreader.Feed
}

func New(feedsFile string) (Scraper, error) {
	scraper := Scraper{}

	feeds, err := loadFeeds(feedsFile)
	if err != nil {
		return scraper, err
	}

	scraper.Feeds = feeds
	return scraper, nil
}

// Scrape downloads the content of the provide list of urls
func (scraper *Scraper) Scrape() {
	concurrencyLimit := 500
	wg := sync.WaitGroup{}
	queue := make(chan *feedreader.FeedItem)

	items := 0
	for _, feed := range scraper.Feeds {
		items += len(feed.Items)
	}
	bar := pb.StartNew(items)

	// prevents "Unsolicited response" log messages from http package when encountering buggy webserver
	log.SetOutput(ioutil.Discard)

	for worker := 0; worker < concurrencyLimit; worker++ {
		wg.Add(1)

		go func(worker int) {
			defer wg.Done()

			for item := range queue {
				err := fetchItem(item)
				bar.Increment()
				if err != nil {
				}
			}
		}(worker)
	}

	for _, feed := range scraper.Feeds {
		for _, item := range feed.Items {
			queue <- item
		}
	}

	close(queue)
	wg.Wait()
	bar.Finish()
	log.SetOutput(os.Stderr)
}

func fetchItem(item *feedreader.FeedItem) error {
	page, err := fetchPage(item.URL)
	if err != nil {
		return err
	}

	content, err := extractContent(item.URL, page)
	if err != nil {
		return err
	}

	item.Content = content

	return nil
}

func (scraper *Scraper) Store(outDir string, location *time.Location) error {
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		err := os.MkdirAll(outDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

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

func fetchPage(url string) (string, error) {
	timeout := time.Duration(20 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return "", fmt.Errorf("Site returned status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func loadFeeds(path string) ([]feedreader.Feed, error) {
	articlesFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var feeds []feedreader.Feed
	err = json.Unmarshal(articlesFile, &feeds)
	if err != nil {
		return nil, err
	}

	return feeds, nil
}
