package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/thesoenke/news-crawler/feedreader"
	"gopkg.in/cheggaaa/pb.v1"
	"gopkg.in/olivere/elastic.v5"
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
func (scraper *Scraper) Scrape(outDir string, location *time.Location, elasticClient *elastic.Client, verbose bool) error {
	wg := sync.WaitGroup{}
	queue := make(chan *feedreader.FeedItem)
	errChan := make(chan bool)
	articleChan := make(chan *Article)

	numItems := 0
	for _, feed := range scraper.Feeds {
		numItems += len(feed.Items)
	}
	bar := pb.StartNew(numItems)

	if !verbose {
		// prevents "Unsolicited response" log messages from http package when encountering buggy webserver
		log.SetOutput(ioutil.Discard)
	}

	startWorker(&wg, queue, articleChan, errChan, verbose)
	go fillWorker(queue, scraper.Feeds)

	failures := 0
	for i := 0; i < numItems; i++ {
		select {
		case article := <-articleChan:
			err := article.Write(outDir, location)
			if err != nil {
				log.SetOutput(os.Stderr)
				return err
			}

			err = article.Index(elasticClient)
			if err != nil {
				log.SetOutput(os.Stderr)
				return err
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

	scraper.Articles = numItems
	scraper.Failures = failures

	return nil
}

func startWorker(wg *sync.WaitGroup, queue chan *feedreader.FeedItem, articleChan chan *Article, errChan chan bool, verbose bool) {
	concurrencyLimit := 100

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

				articleChan <- article
			}
		}(verbose)
	}
}

func fillWorker(queue chan *feedreader.FeedItem, feeds []feedreader.Feed) {
	items := make([]*feedreader.FeedItem, 0)

	for _, feed := range feeds {
		for _, item := range feed.Items {
			items = append(items, item)
		}
	}

	shuffle(items)

	for _, item := range items {
		queue <- item
	}
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

func shuffle(items []*feedreader.FeedItem) {
	for i := range items {
		j := rand.Intn(i + 1)
		items[i], items[j] = items[j], items[i]
	}
}
