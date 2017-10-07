package scraper

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type FetchError struct {
	Msg    string    `json:"message"`
	URL    string    `json:"url"`
	Status int       `json:"status"`
	Time   time.Time `json:"time"`
}

func (e *FetchError) Error() string {
	return fmt.Sprintf("%v %s", e.Time, e.Msg)
}

// Fetch the content of an article from the web
func (article *Article) Fetch() error {
	timeout := time.Duration(60 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest("GET", article.FeedItem.URL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		var status int
		if resp != nil {
			status = resp.StatusCode
		}

		fetchErr := &FetchError{
			Msg:    err.Error(),
			URL:    article.FeedItem.URL,
			Status: status,
			Time:   time.Now(),
		}
		return fetchErr
	}

	defer resp.Body.Close()
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		fetchErr := &FetchError{
			Msg:    fmt.Sprintf("Server returned status code %d", resp.StatusCode),
			URL:    article.FeedItem.URL,
			Status: resp.StatusCode,
			Time:   time.Now(),
		}
		return fetchErr
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	article.HTML = string(body)
	return nil
}
