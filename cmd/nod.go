package cmd

import (
	"path"

	"github.com/spf13/cobra"
	"github.com/thesoenke/news-crawler/nod"
)

var nodOutDir string
var cmdNoD = &cobra.Command{
	Use: "nod",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := path.Join(nodOutDir, lang)
		err := nod.CreateCorpus(lang, dir)
		return err
	},
}

func init() {
	cmdNoD.PersistentFlags().StringVarP(&nodOutDir, "dir", "d", "out/nod", "Directory to store daily compressed text corpus")
	cmdNoD.PersistentFlags().StringVarP(&lang, "lang", "l", "english", "Language of the content")
	RootCmd.AddCommand(cmdNoD)
}
