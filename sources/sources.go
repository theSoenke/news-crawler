package sources

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

func Run() error {
	feedDirectories := fetchFeedDirectories()
	feeds, err := collectFeeds(feedDirectories)
	if err != nil {
		return err
	}

	fmt.Printf("Found %d feeds\n", len(feeds))

	err = store(feeds)
	return err
}

func store(feeds []string) error {
	feedsJSON, err := json.Marshal(feeds)
	if err != nil {
		return err
	}
	ioutil.WriteFile("feeds.json", feedsJSON, 0644)
	return nil
}

func fetchFeedDirectories() []string {
	urls := make([]string, 0)

	// http://www.rss-verzeichnis.net/
	for i := 1; i < 54; i++ {
		url := fmt.Sprintf("http://www.rss-verzeichnis.net/nachrichten-page%d.htm", i)
		urls = append(urls, url)
	}

	return urls
}

func collectFeeds(directories []string) ([]string, error) {
	feedURLs := make([]string, 0)
	for _, url := range directories {
		pageHTML, err := retrievePage(url)
		if err != nil {
			return nil, err
		}
		urls := extractFeeds(pageHTML)
		feedURLs = append(feedURLs, urls...)
	}

	return feedURLs, nil
}

func retrievePage(url string) (string, error) {
	rsp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	html, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}

	return string(html), nil
}

func extractFeeds(html string) []string {
	feedReg := regexp.MustCompile(`
		(https?:\/\/
		([-\w\.]+)+(:\d+)?
		(\/([\w\/_\.]*(\?\S+)?)?)?
		(feed|rss)+
		([\w\/_\.\-]*(\?\S+)?)?)`)
	feeds := feedReg.FindAllString(html, -1)
	feeds = uniq(feeds)

	return feeds
}

func uniq(s []string) []string {
	seen := make(map[string]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
}
