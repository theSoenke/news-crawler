package cmd

import (
	"path"
	"time"

	"github.com/spf13/cobra"
	"github.com/thesoenke/news-crawler/nod"
)

var nodOutDir string
var fromDate string
var cmdNoD = &cobra.Command{
	Use: "nod",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := path.Join(nodOutDir, lang)

		if fromDate == "yesterday" {
			location, err := time.LoadLocation(timeZone)
			if err != nil {
				return err
			}

			yesterday := time.Now().In(location).AddDate(0, 0, -1)
			fromDate = yesterday.Format("2-1-2006")
		}

		err := nod.CreateCorpus(lang, fromDate, timeZone, dir)
		return err
	},
}

func init() {
	cmdNoD.PersistentFlags().StringVarP(&nodOutDir, "dir", "d", "out/nod", "Directory to store daily compressed text corpus")
	cmdNoD.PersistentFlags().StringVarP(&lang, "lang", "l", "", "Language of the content")
	cmdNoD.PersistentFlags().StringVarP(&fromDate, "from", "f", "25-08-2017", "Start date for extraction")
	cmdNoD.PersistentFlags().StringVarP(&timeZone, "timezone", "t", "Europe/Berlin", "Timezone for storing the feeds")
	RootCmd.AddCommand(cmdNoD)
}
