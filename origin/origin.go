package origin

import (
	"github.com/juzeon/lip/data"
	"github.com/juzeon/lip/util"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"net"
	"strings"
)

type Origin struct {
	Name             string
	DatabaseURLs     []string
	DatabaseFileName string
	LookUpFunc       func(ip net.IP) (*data.IPLookupResult, error)
}

var Origins = []Origin{
	{
		Name: "ip2region",
		DatabaseURLs: []string{"https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region.xdb",
			"https://ghproxy.com/github.com/lionsoul2014/ip2region/raw/master/data/ip2region.xdb",
			"https://raw.fastgit.org/lionsoul2014/ip2region/master/data/ip2region.xdb"},
		DatabaseFileName: "ip2region.xdb",
		LookUpFunc: func(ip net.IP) (*data.IPLookupResult, error) {
			searcher, err := xdb.NewWithFileOnly(util.MustLipPath("ip2region.xdb"))
			if err != nil {
				return nil, err
			}
			defer searcher.Close()
			str, err := searcher.SearchByStr(ip.String())
			if err != nil {
				return nil, err
			}
			arr := strings.Split(str, "|")
			return &data.IPLookupResult{
				Country: arr[0],
				State:   arr[2],
				City:    arr[3],
				ISP:     arr[4],
			}, nil
		},
	},
}
