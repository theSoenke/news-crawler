package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/spf13/cobra"
	"github.com/thesoenke/news-crawler/scraper"
)

var scrapeOutDir string
var cmdScrape = &cobra.Command{
	Use:   "scrape",
	Short: "Scrape content from provided URLs",
	RunE: func(cmd *cobra.Command, args []string) error {
		location, err := time.LoadLocation(timeZone)
		if err != nil {
			return err
		}

		yesterday := time.Now().In(location).AddDate(0, 0, -1)
		feedPath, err := getFeedsFilePath(args[0], &yesterday)
		if err != nil {
			return err
		}

		contentScraper, err := scraper.New(feedPath)
		if err != nil {
			return err
		}

		contentScraper.Lang = lang
		contentScraper.Verbose = verbose
		start := time.Now()
		dir := path.Join(scrapeOutDir, lang)
		err = contentScraper.Scrape(dir, &yesterday)
		if err != nil {
			return err
		}

		successLog := fmt.Sprintf("Scraper %s\nArticles: %d\nFailures: %d\nTime: %s\nFile: %s\n", time.Now().In(location), contentScraper.Articles-contentScraper.Failures, contentScraper.Failures, time.Since(start), feedPath)
		fmt.Println(successLog)
		err = writeLog(logsDir, successLog)
		return err
	},
}

func init() {
	cmdScrape.Args = cobra.ExactArgs(1)
	cmdScrape.PersistentFlags().StringVarP(&scrapeOutDir, "dir", "d", "out/content/", "Directory to store fetched pages")
	cmdScrape.PersistentFlags().StringVar(&logsDir, "logs", "out/events.log", "File to store logs")
	cmdScrape.PersistentFlags().StringVarP(&lang, "lang", "l", "", "Language of the content")
	cmdScrape.PersistentFlags().StringVarP(&timeZone, "timezone", "t", "Europe/Berlin", "Timezone for storing the feeds")
	cmdScrape.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose logging of scraper")
	RootCmd.AddCommand(cmdScrape)
}

func getFeedsFilePath(feedsPath string, day *time.Time) (string, error) {
	if feedsPath == "" {
		return "", errors.New("Please provide a file or directory with feed articles")
	}

	stat, err := os.Stat(feedsPath)
	if err != nil {
		return "", err
	}

	// Use article list from yesterday as the input file
	// This ensures that all articles for one day are included
	if stat.IsDir() {
		dayStr := day.Format("2-1-2006")
		path := path.Join(feedsPath, dayStr+".json")
		_, err := os.Stat(path)
		return path, err
	}

	return feedsPath, nil
}
