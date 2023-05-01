package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/juzeon/lip/data"
	"github.com/juzeon/lip/httpclient"
	"github.com/juzeon/lip/source"
	"github.com/juzeon/lip/util"
	"github.com/spf13/cobra"
	"log"
	"net"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "lip <IP or domain>",
	Short: "A tool to look up IP addresses.",
	Long: `lip is a versatile command-line interface (CLI) tool that enables users to 
look up IP addresses and perform a wide range of additional 
functions. With lip, users can easily look up IP addresses, 
both for IPv4 and IPv6, and obtain detailed information about 
the associated domain names, subnets, and geolocations.`,
	Args: cobra.MatchAll(cobra.ExactArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		if persistentFlags.Proxy != "" {
			httpclient.SetProxy(persistentFlags.Proxy)
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
		if rootFlags.Both {
			fmt.Println("Fetching results from both offline and online sources...")
		}
		for _, ip := range ips {
			resArr := doLookup(ip, false)
			if rootFlags.Both && !rootFlags.Offline {
				resArr = append(resArr, doLookup(ip, true)...)
			}
			fmt.Println(util.Ternary(rootFlags.Both, "Lookup", "Offline lookup") +
				" result of " + ip.String() + ": ")
			renderIPLookupResult(resArr)
		}
		if !rootFlags.Both && !rootFlags.Offline {
			fmt.Println("Fetching results from online sources...")
			for _, ip := range ips {
				resArr := doLookup(ip, true)
				fmt.Println("Online lookup result of " + ip.String() + ": ")
				renderIPLookupResult(resArr)
			}
		}
	},
}

type rootFlagStruct struct {
	Reverse bool
	Both    bool
	NoCache bool
	Offline bool
	JSON    bool
}

var rootFlags = rootFlagStruct{}

type persistentFlagStruct struct {
	Proxy string
}

var persistentFlags = persistentFlagStruct{}

func init() {
	rootCmd.PersistentFlags().StringVarP(&persistentFlags.Proxy, "proxy", "p", "",
		"set up a proxy, for example: http://127.0.0.1:7890")
	rootCmd.Flags().BoolVarP(&rootFlags.Reverse, "reverse", "r", false,
		"reverse the output table")
	rootCmd.Flags().BoolVarP(&rootFlags.Both, "both", "b", false,
		"look up an IP or domain from both offline and online sources at once")
	rootCmd.Flags().BoolVarP(&rootFlags.NoCache, "nocache", "n", false,
		"disable cache for the IP to look up (only resolved IPs and online sources will be cached)")
	rootCmd.Flags().BoolVarP(&rootFlags.Offline, "offline", "O", false,
		"look up an IP from offline sources only")
	rootCmd.Flags().BoolVarP(&rootFlags.JSON, "json", "j", false,
		"use JSON output format instead of ASCII table")
}
func renderIPLookupResult(resArr []data.IPLookupResult) {
	if rootFlags.JSON {
		v, err := json.MarshalIndent(resArr, "", "  ")
		if err != nil {
			log.Fatalln("cannot marshal result json: " + err.Error())
		}
		fmt.Println(string(v))
	} else {
		var matrix [][]string
		matrix = append(matrix, data.IPLookupResultTableHeader)
		for _, res := range resArr {
			matrix = append(matrix, []string{res.Source, res.Country, res.Region, res.City, res.ISP, res.Additional})
		}
		util.WriteTable(matrix, os.Stdout, rootFlags.Reverse)
	}
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
		if !rootFlags.NoCache && src.IsOnline() {
			if cacheRes, ok := source.FindCache(ip, src.GetName()); ok {
				resArr = append(resArr, cacheRes)
				continue
			}
		}
		res, err := src.LookUp(ip)
		if err != nil {
			log.Println("failed to look up IP " + ip.String() + " from source " + src.GetName() + ": " + err.Error())
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
