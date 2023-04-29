package source

import (
	"github.com/juzeon/lip/data"
	"github.com/juzeon/lip/httpclient"
	"github.com/juzeon/lip/util"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"net"
	"strings"
)

type IP2Region struct {
	searcher *xdb.Searcher
}

func (i *IP2Region) GetDatabaseFileName() string {
	return "ip2region.xdb"
}

func (i *IP2Region) Close() {
	i.searcher.Close()
}

func (i *IP2Region) Init() error {
	searcher, err := xdb.NewWithFileOnly(util.MustLipPath(i.GetDatabaseFileName()))
	if err != nil {
		return err
	}
	i.searcher = searcher
	return nil
}

func (i *IP2Region) GetName() string {
	return "ip2region"
}

func (i *IP2Region) DownloadDatabase() error {
	return httpclient.DownloadToWithMultipleURLs([]string{
		"https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region.xdb",
		"https://ghproxy.com/github.com/lionsoul2014/ip2region/raw/master/data/ip2region.xdb",
		"https://raw.fastgit.org/lionsoul2014/ip2region/master/data/ip2region.xdb",
	}, util.MustLipPath(i.GetDatabaseFileName()))
}

func (i *IP2Region) LookUp(ip net.IP, duplicateIdentical bool) (data.IPLookupResult, error) {
	str, err := i.searcher.SearchByStr(ip.String())
	if err != nil {
		return data.IPLookupResult{}, err
	}
	arr := strings.Split(str, "|")
	return data.IPLookupResult{
		Source:  i.GetName(),
		Country: util.Ternary(arr[0] == "0", "", arr[0]),
		State:   util.Ternary(arr[2] == "0", "", arr[2]),
		City:    util.Ternary(arr[3] == "0", "", arr[3]),
		ISP:     util.Ternary(arr[4] == "0", "", arr[4]),
	}, nil
}

func (i *IP2Region) CheckUpdate() bool {
	//TODO implement me
	panic("implement me")
}
