package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var cmdScrape = &cobra.Command{
	Use:   "scrape",
	Short: "Scrape all provided articles",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Print: " + strings.Join(args, " "))

	},
}

func init() {
	RootCmd.AddCommand(cmdScrape)
}
