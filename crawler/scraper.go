package crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// ScrapeURLs downloads the content of the provide list of urls
func ScrapeURLs(urls []string) error {
	for _, url := range urls {
		response, err := http.Get(url)
		if err != nil {
			return err
		}

		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", body)
	}
	return nil
}
