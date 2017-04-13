package crawler

import (
	"log"

	"github.com/mmcdole/gofeed"
)

// Article is a single feed item
type Article struct {
	Title     string
	Content   string
	URL       string
	Published string
}

// Run a new crawler
func Run() {
	sources, err := readSourcesFile("feeds/news_de.json")

	if err != nil {
		log.Fatal("Failed to import sources")
		return
	}

	for _, url := range sources {
		var articles []*Article
		articles, err = parse(url)
		storeArticles(articles)
	}
}

func parse(url string) ([]*Article, error) {
	feedParser := gofeed.NewParser()
	feed, err := feedParser.ParseURL(url)

	if err != nil {
		return nil, err
	}

	articles := make([]*Article, len(feed.Items))
	for i, item := range feed.Items {
		article := Article{
			Title:     item.Title,
			Content:   item.Description,
			URL:       item.Link,
			Published: item.Published,
		}

		articles[i] = &article
	}

	return articles, nil
}

func storeArticles(articles []*Article) {
	// TODO
}
