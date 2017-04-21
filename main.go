package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/thesoenke/news-crawler/feeds/crawler"
)

func main() {
	sources, err := readSourcesFile("rss/news_de.json")
	if err != nil {
		log.Fatal("Failed to import sources")
	}

	feeds, err := crawler.Run(sources)
	if err != nil {
		panic(err)
	}
	saveAsJSON(feeds)
	saveInDB(feeds)
}

func saveAsJSON(feeds []crawler.Feed) {
	for _, feed := range feeds {
		crawler.SaveFeed(feed)
	}
}

func saveInDB(feeds []crawler.Feed) {
	for _, feed := range feeds {
		crawler.Save(feed.Items)
	}
}

func readSourcesFile(path string) ([]string, error) {
	sourceFile, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var sources = make([]string, 0)
	err = json.Unmarshal(sourceFile, &sources)

	if err != nil {
		return nil, err
	}

	return sources, nil
}
