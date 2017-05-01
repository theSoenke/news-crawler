package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"time"

	"github.com/thesoenke/news-crawler/crawler"
)

func main() {
	feedsPtr := flag.String("feeds", "", "Path to a JSON file with a list of feeds")
	timeZonePtr := flag.String("timezone", "", "Timezone for the list of feeds")
	outDirPtr := flag.String("out", "out/", "Directory where to store the output")
	flag.Parse()

	if *feedsPtr == "" {
		log.Fatal("Please provide feed sources with --feeds")
	}

	location, err := time.LoadLocation(*timeZonePtr)
	if err != nil {
		log.Fatal("Please provide a valid timezone with --timezone")
	}

	if *outDirPtr == "" {
		log.Fatal("Please provide a directory to store the output with --out")
	}

	sources, err := readSourcesFile(*feedsPtr)
	if err != nil {
		log.Fatal("Failed to import sources")
	}

	err = crawler.ScrapeFeeds(sources, *outDirPtr, location)
	if err != nil {
		log.Fatal(err.Error())
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
