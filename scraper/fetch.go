package scraper

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	goose "github.com/advancedlogic/GoOse"
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
		fetchErr := &FetchError{
			Msg:    err.Error(),
			URL:    article.FeedItem.URL,
			Status: resp.StatusCode,
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

	contentType := resp.Header.Get("Content-Type")
	charset := getCharsetFromContentType(contentType)
	bodyUtf8 := goose.UTF8encode(string(body), charset)
	article.HTML = bodyUtf8
	return nil
}

func getCharsetFromContentType(cs string) string {
	cs = strings.ToLower(strings.Replace(cs, " ", "", -1))
	if strings.HasPrefix(cs, "text/html;charset=") {
		cs = strings.TrimPrefix(cs, "text/html;charset=")
	}
	if strings.HasPrefix(cs, "text/xhtml;charset=") {
		cs = strings.TrimPrefix(cs, "text/xhtml;charset=")
	}
	if strings.HasPrefix(cs, "application/xhtml+xml;charset=") {
		cs = strings.TrimPrefix(cs, "application/xhtml+xml;charset=")
	}
	if strings.HasPrefix(cs, "text/html") {
		cs = "utf-8"
	}

	return goose.NormaliseCharset(cs)
}
