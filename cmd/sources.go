package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thesoenke/news-crawler/sources"
)

var cmdSources = &cobra.Command{
	Use:   "sources",
	Short: "Scrape feed directories for feeds",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := sources.Run()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	cmdSources.PersistentFlags().StringVarP(&feedsFile, "sources", "s", "", "")
	RootCmd.AddCommand(cmdSources)
}
