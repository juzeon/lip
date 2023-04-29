package source

import (
	"errors"
	"github.com/juzeon/lip/data"
	"github.com/juzeon/lip/httpclient"
	"github.com/juzeon/lip/util"
	"net"
	"strconv"
)

type IPApi struct {
}

func (i *IPApi) Init() error {
	return nil
}

func (i *IPApi) Close() {

}

func (i *IPApi) GetName() string {
	return "ip-api"
}

func (i *IPApi) GetDatabaseFileName() string {
	return ""
}

func (i *IPApi) DownloadDatabase() error {
	return nil
}

func (i *IPApi) LookUp(ip net.IP) (data.IPLookupResult, error) {
	type Result struct {
		Query      string `json:"query"`
		Status     string `json:"status"`
		Message    string `json:"message"`
		Country    string `json:"country"`
		RegionName string `json:"regionName"`
		City       string `json:"city"`
		District   string `json:"district"`
		Isp        string `json:"isp"`
		Org        string `json:"org"`
		As         string `json:"as"`
	}
	resp, err := httpclient.Client.R().SetResult(&Result{}).Get("http://ip-api.com/json/" + ip.String() +
		"?fields=status,message,country,regionName,city,district,isp,org,as,query")
	if err != nil {
		return data.IPLookupResult{}, err
	}
	if resp.StatusCode() != 200 {
		return data.IPLookupResult{}, errors.New("respond with status code: " + strconv.Itoa(resp.StatusCode()))
	}
	result := resp.Result().(*Result)
	if result.Status != "success" {
		return data.IPLookupResult{}, errors.New(result.Message)
	}
	return data.IPLookupResult{
		IP:      ip,
		Source:  i.GetName(),
		Country: result.Country,
		Region:  result.RegionName,
		City:    result.City + util.Ternary(result.District == "", "", ", "+result.District),
		ISP: result.Isp + util.Ternary(result.Org == "", "", ", "+result.Org) +
			util.Ternary(result.As == "", "", ", "+result.As),
	}, nil
}

func (i *IPApi) CheckUpdate() bool {
	return false
}

func (i *IPApi) IsOnline() bool {
	return true
}
