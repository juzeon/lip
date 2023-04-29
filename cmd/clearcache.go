package cmd

import (
	"fmt"
	"github.com/juzeon/lip/source"

	"github.com/spf13/cobra"
)

var clearcacheCmd = &cobra.Command{
	Use:   "clearcache",
	Short: "Clear the cache file of lip",
	Run: func(cmd *cobra.Command, args []string) {
		source.ClearCache()
		fmt.Println("cleared cache file successfully")
	},
}
