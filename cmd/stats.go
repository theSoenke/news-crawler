package cmd

import (
	"github.com/thesoenke/news-crawler/stats"

	"github.com/spf13/cobra"
)

var cmdStats = &cobra.Command{
	Use:   "stats",
	Short: "Return stats about the article index",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := stats.TotalArticles()
		return err
	},
}

func init() {
	RootCmd.AddCommand(cmdStats)
}
