package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"github.com/thesoenke/news-crawler/crawler"
)

func main() {
	feedsPtr := flag.String("feeds", "", "Path to a JSON file with a list of feeds")
	flag.Parse()
	if *feedsPtr == "" {
		log.Fatal("Please provide feed sources with --feeds")
	}

	sources, err := readSourcesFile(*feedsPtr)
	if err != nil {
		log.Fatal("Failed to import sources")
	}

	err = crawler.ScrapeFeeds(sources)
	if err != nil {
		panic(err)
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
