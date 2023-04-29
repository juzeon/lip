package cmd

import (
	"fmt"
	"github.com/juzeon/lip/data"
	"github.com/juzeon/lip/httpclient"
	"github.com/juzeon/lip/source"
	"github.com/juzeon/lip/util"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"log"
	"net"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "lip [IP or domain]",
	Short: "A tool to look up IP addresses",
	Long:  `lip is a tool for looking up IP addresses with many additional functions`,
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		if flags.Proxy != "" {
			httpclient.SetProxy(flags.Proxy)
		}
		util.InitFs()
		source.InitCache()
		source.DownloadDatabases(false)
		source.InitDatabases()
		defer source.CloseDatabases()
		ips, err := parseIP(args[0])
		if err != nil {
			log.Fatalln("cannot parse ip address from the argument: " + err.Error())
		}
		if flags.Both {
			fmt.Println("Fetching results from both offline and online sources...")
		}
		for _, ip := range ips {
			resArr := doLookup(ip, false)
			if flags.Both {
				resArr = append(resArr, doLookup(ip, true)...)
			}
			fmt.Println(util.Ternary(flags.Both, "Lookup", "Offline lookup") +
				" result of " + ip.String() + ": ")
			renderIPLookupResultTable(resArr)
		}
		if !flags.Both {
			fmt.Println("Fetching results from online sources...")
			for _, ip := range ips {
				resArr := doLookup(ip, true)
				fmt.Println("Online lookup result of " + ip.String() + ": ")
				renderIPLookupResultTable(resArr)
			}
		}
	},
}

type flagStruct struct {
	Proxy   string
	Reverse bool
	Both    bool
	NoCache bool
}

var flags = flagStruct{}

func init() {
	rootCmd.PersistentFlags().StringVarP(&flags.Proxy, "proxy", "p", "",
		"set up a proxy, for example: http://127.0.0.1:7890")
	rootCmd.Flags().BoolVarP(&flags.Reverse, "reverse", "r", false,
		"reverse the output table")
	rootCmd.Flags().BoolVarP(&flags.Both, "both", "b", false,
		"look up an IP or domain from both offline and online sources at once")
	rootCmd.Flags().BoolVarP(&flags.NoCache, "nocache", "n", false,
		"disable cache for the IP to look up (only resolved IP and online sources will be cached)")
}
func renderIPLookupResultTable(resArr []data.IPLookupResult) {
	table := tablewriter.NewWriter(os.Stdout)
	var matrix [][]string
	matrix = append(matrix, data.IPLookupResultTableHeader)
	for _, res := range resArr {
		matrix = append(matrix, []string{res.Source, res.Country, res.Region, res.City, res.ISP, res.Additional})
	}
	if flags.Reverse {
		matrix = util.TransposeMatrix(matrix)
	}
	table.AppendBulk(matrix)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_CENTER})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.Render()
}
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
func doLookup(ip net.IP, onlineSource bool) []data.IPLookupResult {
	var resArr []data.IPLookupResult
	for _, src := range source.Sources {
		if onlineSource && !src.IsOnline() ||
			!onlineSource && src.IsOnline() {
			continue
		}
		if !flags.NoCache && src.IsOnline() {
			if cacheRes, ok := source.FindCache(ip, src.GetName()); ok {
				resArr = append(resArr, cacheRes)
				continue
			}
		}
		res, err := src.LookUp(ip)
		if err != nil {
			log.Println("failed to look up IP " + ip.String() + " from source " + src.GetName())
			continue
		}
		resArr = append(resArr, res)
		if src.IsOnline() {
			source.UpsertCache(res)
		}
	}
	return resArr
}
func parseIP(str string) ([]net.IP, error) {
	ip := net.ParseIP(str)
	if ip == nil {
		ips, err := net.LookupIP(str)
		if err != nil {
			return nil, fmt.Errorf("cannot look up ip from domain: %v", err)
		}
		return ips, nil
	}
	return []net.IP{ip}, nil
}
