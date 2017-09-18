package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thesoenke/news-crawler/nod"
)

var cmdNoD = &cobra.Command{
	Use: "nod",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := nod.CreateCorpus("german", "out/nod")
		return err
	},
}

func init() {
	RootCmd.AddCommand(cmdNoD)
}
