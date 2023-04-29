package source

import (
	"github.com/juzeon/lip/data"
	"github.com/juzeon/lip/httpclient"
	"github.com/juzeon/lip/pkg/qqwry"
	"github.com/juzeon/lip/util"
	"net"
)

type QQWry struct {
	db *qqwry.DB
}

func (q *QQWry) Init() error {
	db, err := qqwry.Open(util.MustLipPath(q.GetDatabaseFileName()))
	if err != nil {
		return err
	}
	q.db = db
	return nil
}

func (q *QQWry) Close() {
}

func (q *QQWry) GetName() string {
	return "qqwry"
}

func (q *QQWry) GetDatabaseFileName() string {
	return "qqwry.dat"
}

func (q *QQWry) DownloadDatabase() error {
	return httpclient.DownloadToWithMultipleURLs([]string{
		"https://github.com/FW27623/qqwry/raw/main/qqwry.dat",
		"https://ghproxy.com/github.com/FW27623/qqwry/raw/main/qqwry.dat",
		"https://raw.fastgit.org/FW27623/qqwry/main/qqwry.dat",
	}, util.MustLipPath(q.GetDatabaseFileName()))
}

func (q *QQWry) LookUp(ip net.IP, duplicateIdentical bool) (data.IPLookupResult, error) {
	record, err := q.db.Query(ip)
	if err != nil {
		return data.IPLookupResult{}, err
	}
	result := record.Country + " " + record.City
	return data.IPLookupResult{
		Source:  q.GetName(),
		Country: result,
		State:   util.Ternary(duplicateIdentical, result, ""),
		City:    util.Ternary(duplicateIdentical, result, ""),
		ISP:     util.Ternary(duplicateIdentical, result, ""),
	}, nil
}

func (q *QQWry) CheckUpdate() bool {
	//TODO implement me
	panic("implement me")
}
