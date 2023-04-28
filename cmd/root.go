package cmd

import (
	"fmt"
	"github.com/juzeon/lip/httpclient"
	"github.com/juzeon/lip/storage"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lip [IP or domain]",
	Short: "A tool to look up IP addresses",
	Long:  `lip is a tool for looking up IP addresses with many additional functions`,
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		if Proxy != "" {
			httpclient.SetProxy(Proxy)
		}
		storage.DownloadDatabases(false)

	},
}

var Proxy string

func init() {
	rootCmd.PersistentFlags().StringVarP(&Proxy, "proxy", "p", "", "设置代理URL，如：http://127.0.0.1:7890")
}
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
