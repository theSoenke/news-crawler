package sources

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

func Run(url string) error {
	pageHTML, err := retrieve(url)
	if err != nil {
		return err
	}

	extractFeeds(pageHTML)
	return nil
}

func retrieve(url string) (string, error) {
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
	feedReg := regexp.MustCompile(`(https?://([-\w\.]+)+(:\d+)?(/([\w/_\.]*(\?\S+)?)?)?/(feed|rss))`)
	feedLinks := feedReg.FindAllString(html, -1)
	feedLinks = feedsUniq(feedLinks)
	fmt.Println(feedLinks)

	// TODO
	return feedLinks
}

func feedsUniq(s []string) []string {
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
