package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

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
		return scraper, nil
	}

	scraper.Feeds = feeds
	return scraper, nil
}

// Scrape downloads the content of the provide list of urls
func (scraper *Scraper) Scrape() error {
	items := 0
	for _, feed := range scraper.Feeds {
		items += len(feed.Items)
	}
	bar := pb.StartNew(items)

	for _, feed := range scraper.Feeds {
		for _, feedItem := range feed.Items {
			page, err := fetchPage(feedItem.URL)
			if err != nil {
				return err
			}

			content, err := extractContent(feedItem.URL, page)
			if err != nil {
				return err
			}

			feedItem.Content = content
			bar.Increment()
		}
	}

	bar.Finish()
	store(scraper.Feeds)

	return nil
}

func store(feeds []feedreader.Feed) error {
	feedsJSON, err := json.Marshal(feeds)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("content.json", feedsJSON, 0644)
	if err != nil {
		return err
	}
	return nil
}

func fetchPage(url string) (string, error) {
	timeout := time.Duration(30 * time.Second)
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
		return "", fmt.Errorf("Site retuned status code %d", resp.StatusCode)
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
