package crawler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/mmcdole/gofeed"
)

type article struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	URL       string `json:"url"`
	Published string `json:"published"`
}

// Run a new crawler
func Run() {
	sources, err := readSourcesFile("feeds/news_de.json")

	if err != nil {
		log.Fatal("Failed to import sources")
	}

	for _, url := range sources {
		var articles []*article
		articles, err = parse(url)
		err = save(url, articles)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(url)
	}
}

func parse(url string) ([]*article, error) {
	feedParser := gofeed.NewParser()
	feed, err := feedParser.ParseURL(url)

	if err != nil {
		return nil, err
	}

	articles := make([]*article, len(feed.Items))
	for i, item := range feed.Items {
		newArticle := article{
			Title:     item.Title,
			Content:   item.Description,
			URL:       item.Link,
			Published: item.Published,
		}

		articles[i] = &newArticle
	}

	return articles, nil
}

func save(feedURL string, articles []*article) error {
	jsonArticles, err := json.Marshal(articles)
	if err != nil {
		return err
	}

	path := filepath.Join(".", "out")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}

	var domain *url.URL
	domain, err = url.Parse(feedURL)

	outFilePath := "out/" + domain.Host + ".json"
	err = ioutil.WriteFile(outFilePath, jsonArticles, 0644)
	return err
}
