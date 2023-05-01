package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/juzeon/lip/util"
	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var whoisCmd = &cobra.Command{
	Use:   "whois [domain]",
	Short: "WHOIS lookup for a domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dialer, err := util.GetProxyDialer(persistentFlags.Proxy,
			time.Duration(whoisFlags.Timeout)*time.Second)
		if err != nil {
			log.Fatalln("cannot get dialer: " + err.Error())
		}
		client := whois.NewClient().SetDialer(dialer).SetTimeout(time.Duration(whoisFlags.Timeout) * time.Second)
		rawResult, err := client.Whois(args[0])
		if err != nil {
			log.Fatalln("cannot do whois lookup: " + err.Error())
		}
		info, err := whoisparser.Parse(rawResult)
		if err != nil {
			log.Println("cannot parse whois result. Use raw result instead: " + err.Error())
			fmt.Println(rawResult)
		} else {
			renderWhoisResult(info)
		}
	},
}

func renderWhoisResult(info whoisparser.WhoisInfo) {
	if whoisFlags.JSON {
		var target any
		if whoisFlags.Domain {
			if info.Domain == nil {
				log.Println("no information for domain")
				return
			}
			target = info.Domain
		} else {
			target = info
		}
		v, err := json.MarshalIndent(target, "", "  ")
		if err != nil {
			log.Fatalln("cannot marshal whois json: " + err.Error())
		}
		fmt.Println(string(v))
		return
	}
	if info.Domain != nil {
		fmt.Println("Domain:")
		var data [][]string
		data = append(data, removeWhoisRedacted([]string{"" +
			whoisRedact("ID"),
			"Domain",
			"Punycode",
			whoisRedact("Name"),
			whoisRedact("Extension"),
			"WhoisServer",
			"Status",
			"NameServers",
			"DNSSec",
			"CreatedDate",
			"UpdatedDate",
			"ExpirationDate",
		}))
		data = append(data, removeWhoisRedacted([]string{
			whoisRedact(info.Domain.ID),
			info.Domain.Domain,
			info.Domain.Punycode,
			whoisRedact(info.Domain.Name),
			whoisRedact(info.Domain.Extension),
			info.Domain.WhoisServer,
			strings.Join(info.Domain.Status, ", "),
			strings.Join(info.Domain.NameServers, ", "),
			util.Ternary(info.Domain.DNSSec, "true", "false"),
			info.Domain.CreatedDate,
			info.Domain.UpdatedDate,
			info.Domain.ExpirationDate,
		}))
		util.WriteTable(data, os.Stdout, whoisFlags.Reverse)
	}
	if whoisFlags.Domain {
		return
	}
	type ContactWithName struct {
		Name    string
		Contact *whoisparser.Contact
	}
	var data [][]string
	data = append(data, removeWhoisRedacted([]string{
		"Source",
		whoisRedact("ID"),
		"Name",
		"Organization",
		"Street",
		"City",
		"Province",
		"PostalCode",
		"Country",
		"Phone",
		whoisRedact("PhoneExt"),
		whoisRedact("Fax"),
		whoisRedact("FaxExt"),
		"Email",
		"ReferralURL",
	}))
	for _, contact := range []ContactWithName{
		{Contact: info.Registrar, Name: "Registrar"}, {Contact: info.Registrant, Name: "Registrant"},
		{Contact: info.Administrative, Name: "Administrative"}, {Contact: info.Technical, Name: "Technical"},
		{Contact: info.Billing, Name: "Billing"}} {
		if contact.Contact == nil {
			continue
		}
		data = append(data, removeWhoisRedacted([]string{
			contact.Name,
			whoisRedact(contact.Contact.ID),
			contact.Contact.Name,
			contact.Contact.Organization,
			contact.Contact.Street,
			contact.Contact.City,
			contact.Contact.Province,
			contact.Contact.PostalCode,
			contact.Contact.Country,
			contact.Contact.Phone,
			whoisRedact(contact.Contact.PhoneExt),
			whoisRedact(contact.Contact.Fax),
			whoisRedact(contact.Contact.FaxExt),
			contact.Contact.Email,
			contact.Contact.ReferralURL,
		}))
	}
	if len(data) > 1 {
		if info.Domain != nil {
			fmt.Println()
		}
		fmt.Println("Contacts:")
	}
	util.WriteTable(data, os.Stdout, whoisFlags.Reverse)
}
func whoisRedact(str string) string {
	return util.Ternary(whoisFlags.Full, str, "^")
}
func removeWhoisRedacted(arr []string) []string {
	return util.FilterKeep(arr, func(value string) bool {
		return value != "^"
	})
}

type whoisFlagStruct struct {
	Reverse bool
	Timeout int
	JSON    bool
	Full    bool
	Domain  bool
}

var whoisFlags = whoisFlagStruct{}

func init() {
	rootCmd.AddCommand(whoisCmd)

	whoisCmd.Flags().BoolVarP(&whoisFlags.Reverse, "reverse", "r", false,
		"reverse the output table")
	whoisCmd.Flags().IntVarP(&whoisFlags.Timeout, "timeout", "t", 10,
		"timeout for whois lookup")
	whoisCmd.Flags().BoolVarP(&whoisFlags.JSON, "json", "j", false,
		"use JSON output format instead of ASCII table")
	whoisCmd.Flags().BoolVarP(&whoisFlags.Full, "full", "f", false,
		"show full information columns of whois lookup")
	whoisCmd.Flags().BoolVarP(&whoisFlags.Domain, "domain", "d", false,
		"show information of the domain only, hiding contacts")
}
