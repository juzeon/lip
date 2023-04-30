package source

import (
	"errors"
	"github.com/juzeon/lip/data"
	"github.com/juzeon/lip/httpclient"
	"github.com/juzeon/lip/util"
	"net"
	"strconv"
	"strings"
	"time"
)

type IPInfo struct {
}

func (i *IPInfo) Init() error {
	return nil
}

func (i *IPInfo) Close() {
}

func (i *IPInfo) GetName() string {
	return "ipinfo"
}

func (i *IPInfo) GetDatabaseFileName() string {
	return ""
}

func (i *IPInfo) DownloadDatabase() error {
	return nil
}

func (i *IPInfo) LookUp(ip net.IP) (data.IPLookupResult, error) {
	resp, err := httpclient.GetClient().SetTimeout(10*time.Second).R().
		SetResult(&ipinfoResult{}).
		SetHeader("Referer", "https://ipinfo.io/").
		Get("https://ipinfo.io/widget/demo/" + ip.String())
	if err != nil {
		return data.IPLookupResult{}, err
	}
	if resp.StatusCode() != 200 {
		return data.IPLookupResult{}, errors.New("respond with status code: " + strconv.Itoa(resp.StatusCode()))
	}
	result := resp.Result().(*ipinfoResult)
	return data.IPLookupResult{
		IP:      ip,
		Source:  i.GetName(),
		Country: result.Data.Country,
		Region:  result.Data.Region,
		City:    result.Data.City,
		ISP: util.Ternary(result.Data.Asn.Asn == "", "", " "+result.Data.Asn.Asn) +
			util.Ternary(result.Data.Asn.Name == "", "", " "+result.Data.Asn.Name) +
			util.Ternary(result.Data.Asn.Domain == "", "", " "+result.Data.Asn.Domain) +
			util.Ternary(result.Data.Company.Name == "", "", ", "+result.Data.Company.Name),
		Additional: util.Ternary(result.Data.Hostname == "", "", "hostname: "+result.Data.Hostname+"\n") +
			util.Ternary(len(result.Data.Domains.Domains) == 0, "", "domains: "+
				strings.Join(result.Data.Domains.Domains, ", ")+"\n"),
	}, nil
}

func (i *IPInfo) CheckUpdate() bool {
	return false
}

func (i *IPInfo) IsOnline() bool {
	return true
}

type ipinfoResult struct {
	Input string `json:"input"`
	Data  struct {
		Ip       string `json:"ip"`
		Hostname string `json:"hostname"`
		City     string `json:"city"`
		Region   string `json:"region"`
		Country  string `json:"country"`
		Loc      string `json:"loc"`
		Org      string `json:"org"`
		Postal   string `json:"postal"`
		Timezone string `json:"timezone"`
		Asn      struct {
			Asn    string `json:"asn"`
			Name   string `json:"name"`
			Domain string `json:"domain"`
			Route  string `json:"route"`
			Type   string `json:"type"`
		} `json:"asn"`
		Company struct {
			Name   string `json:"name"`
			Domain string `json:"domain"`
			Type   string `json:"type"`
		} `json:"company"`
		Privacy struct {
			Vpn     bool   `json:"vpn"`
			Proxy   bool   `json:"proxy"`
			Tor     bool   `json:"tor"`
			Relay   bool   `json:"relay"`
			Hosting bool   `json:"hosting"`
			Service string `json:"service"`
		} `json:"privacy"`
		Domains struct {
			Total   int      `json:"total"`
			Domains []string `json:"domains"`
		} `json:"domains"`
	} `json:"data"`
}
