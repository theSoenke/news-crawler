package crawler

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Save will store a feed item in the database
func Save(items []*FeedItem) {
	db, err := sql.Open("sqlite3", "./feeds.db")
	checkErr(err)

	createTables(db)
	checkErr(err)

	for _, item := range items {
		saveItem(db, item)
	}
}

func saveItem(db *sql.DB, item *FeedItem) {
	stmt, err := db.Prepare("INSERT INTO feeds(guid, url, created) values(?,?,?)")
	checkErr(err)

	created := time.Now().Format("01-02-2006")
	_, err = stmt.Exec(item.GUID, item.URL, created)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			fmt.Println(item.GUID)
			fmt.Println(err.Error())
		} else {
			checkErr(err)
		}
	}
}

func createTables(db *sql.DB) error {
	sqlTable := `
	CREATE TABLE feeds (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		guid TEXT NOT NULL UNIQUE,
		url TEXT,
		created DATETIME
	);
	`

	_, err := db.Exec(sqlTable)
	return err
}

// SaveFeed will save the content of a feed as JSON
func SaveFeed(feed Feed) error {
	feedItems, err := json.Marshal(feed.Items)
	if err != nil {
		return err
	}

	path := filepath.Join(".", "out")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}

	var domain *url.URL
	domain, err = url.Parse(feed.URL)

	outFilePath := "out/" + domain.Host + ".json"
	err = ioutil.WriteFile(outFilePath, feedItems, 0644)
	return err
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
