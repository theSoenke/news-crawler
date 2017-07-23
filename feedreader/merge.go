package feedreader

import (
	"bufio"
	"os"
)

func loadFeeds(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	feeds := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		feeds = append(feeds, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return feeds, nil
}

func merge(newFeeds []Feed, oldFeeds []Feed) []Feed {
	oldFeedsMap := make(map[string][]*FeedItem, len(oldFeeds))
	for _, feed := range oldFeeds {
		oldFeedsMap[feed.URL] = feed.Items
	}

	feeds := make([]Feed, 0)
	for _, feed := range newFeeds {
		if _, ok := oldFeedsMap[feed.URL]; ok {
			items := removeDuplicates(feed.Items, oldFeedsMap[feed.URL])
			newFeed := Feed{
				URL:   feed.URL,
				Items: items,
			}
			feeds = append(feeds, newFeed)
		} else {
			feeds = append(feeds, feed)
		}
	}

	return feeds
}

func removeDuplicates(newItems []*FeedItem, oldItems []*FeedItem) []*FeedItem {
	found := make(map[string]bool)

	for _, item := range oldItems {
		found[item.GUID] = true
	}

	for _, item := range newItems {
		if !found[item.GUID] {
			oldItems = append(oldItems, item)
		}
	}

	return oldItems
}
