package scraper

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"fmt"

	"github.com/thesoenke/news-crawler/feedreader"
	"gopkg.in/cheggaaa/pb.v1"
)

const (
	userAgent = "Mozilla/5.0 (X11; Linux x86_64; rv:53.0) Gecko/20100101 Firefox/53.0"
)

type Scraper struct {
	Feeds    []feedreader.Feed
	Articles int
	Failures int
}

type Article struct {
	FeedItem *feedreader.FeedItem
	HTML     string
}

// New creates a scraper instance
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
func (scraper *Scraper) Scrape(verbose bool) error {
	concurrencyLimit := 500
	wg := sync.WaitGroup{}
	queue := make(chan *feedreader.FeedItem)
	errChan := make(chan bool)
	ch := make(chan *Article)

	articles := 0
	for _, feed := range scraper.Feeds {
		articles += len(feed.Items)
	}
	bar := pb.StartNew(articles)

	if !verbose {
		// prevents "Unsolicited response" log messages from http package when encountering buggy webserver
		log.SetOutput(ioutil.Discard)
	}

	for worker := 0; worker < concurrencyLimit; worker++ {
		wg.Add(1)

		go func(verbose bool) {
			defer wg.Done()

			for item := range queue {
				article := &Article{
					FeedItem: item,
				}

				err := article.Fetch()
				if err != nil {
					if verbose {
						fmt.Printf("Failed to fetch %s %s\n", item.URL, err)
					}
					errChan <- true
					continue
				}

				err = article.Extract()
				if err != nil {
					if verbose {
						fmt.Printf("Failed to extract %s %s\n", item.URL, err)
					}
					errChan <- true
					continue
				}

				ch <- article
			}
		}(verbose)
	}

	go func(feeds []feedreader.Feed) {
		for _, feed := range feeds {
			for _, item := range feed.Items {
				queue <- item
			}
		}
	}(scraper.Feeds)

	elasticClient, err := NewElasticClient()
	if err != nil {
		log.SetOutput(os.Stderr)
		return err
	}

	failures := 0
	for i := 0; i < articles; i++ {
		select {
		case article := <-ch:
			err := article.StoreElastic(elasticClient)
			if err != nil {
				log.SetOutput(os.Stderr)
				log.Fatal(err)
			}
		case <-errChan:
			failures++
		}
		bar.Increment()
	}

	close(queue)
	wg.Wait()
	bar.Finish()
	log.SetOutput(os.Stderr)

	scraper.Articles = articles
	scraper.Failures = failures

	return nil
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
