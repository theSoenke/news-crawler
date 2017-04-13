package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	// _, err := Parse("http://spiegel.de")

	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	sources, err := readSourcesFile("feeds/news_de.json")

	if err != nil {
		log.Fatal("Failed to import sources")
		return
	}

	fmt.Print(sources)
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
