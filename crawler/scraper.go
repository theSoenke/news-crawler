package crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	userAgent = "Mozilla/5.0 (X11; Linux x86_64; rv:53.0) Gecko/20100101 Firefox/53.0"
)

// ScrapeURLs downloads the content of the provide list of urls
func ScrapeURLs(urls []string) error {
	for _, url := range urls {
		client := http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", userAgent)
		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", body)
	}
	return nil
}
