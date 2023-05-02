package source

import (
	"errors"
	"fmt"
	"github.com/ipipdotnet/ipdb-go"
	"github.com/juzeon/lip/data"
	"github.com/juzeon/lip/util"
	"net"
)

type IPIP struct {
	city *ipdb.City
}

func (i *IPIP) Init() error {
	if !util.MustCheckExist(util.MustLipPath(i.GetDatabaseFileName())) {
		return errors.New("please manually download ipip database to " + util.MustLipPath(i.GetDatabaseFileName()) +
			" from https://www.ipip.net/product/ip.html")
	}
	city, err := ipdb.NewCity(util.MustLipPath(i.GetDatabaseFileName()))
	if err != nil {
		return fmt.Errorf("cannot init ipip: " + err.Error())
	}
	i.city = city
	return nil
}

func (i *IPIP) Close() {
}

func (i *IPIP) GetName() string {
	return "ipip"
}

func (i *IPIP) GetDatabaseFileName() string {
	return "ipip.ipdb"
}

func (i *IPIP) DownloadDatabase() error {
	return nil
}

func (i *IPIP) LookUp(ip net.IP) (data.IPLookupResult, error) {
	info, err := i.city.FindInfo(ip.String(), "CN")
	if err != nil {
		return data.IPLookupResult{}, err
	}
	return data.IPLookupResult{
		IP:         ip,
		Source:     i.GetName(),
		Country:    info.CountryName,
		Region:     info.RegionName,
		City:       info.CityName,
		ISP:        "",
		Additional: "",
	}, nil
}

func (i *IPIP) CheckUpdate() bool {
	return false
}

func (i *IPIP) IsOnline() bool {
	return false
}
