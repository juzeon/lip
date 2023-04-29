package source

import (
	"encoding/json"
	"github.com/juzeon/lip/data"
	"github.com/juzeon/lip/util"
	"log"
	"net"
	"os"
)

var cacheData []data.IPLookupResult

func InitCache() {
	if util.MustCheckExist(util.MustLipPath("cache.json")) {
		v, err := os.ReadFile(util.MustLipPath("cache.json"))
		if err != nil {
			log.Fatalln("cannot read cache file: " + err.Error())
		}
		err = json.Unmarshal(v, &cacheData)
		if err != nil {
			log.Fatalln("cannot unmarshal cache file: " + err.Error())
		}
	}
}
func FindCache(ip net.IP, source string) (data.IPLookupResult, bool) {
	for _, res := range cacheData {
		if res.IP.String() == ip.String() && res.Source == source {
			return res, true
		}
	}
	return data.IPLookupResult{}, false
}
func UpsertCache(res data.IPLookupResult) {
	updated := false
	for ix, cacheRes := range cacheData {
		if res.IP.String() == cacheRes.IP.String() && res.Source == cacheRes.Source {
			cacheData[ix] = res
			updated = true
		}
	}
	if !updated {
		cacheData = append(cacheData, res)
	}
	v, err := json.MarshalIndent(&cacheData, "", "  ")
	if err != nil {
		log.Fatalln("cannot marshal cache data: " + err.Error())
	}
	err = os.WriteFile(util.MustLipPath("cache.json"), v, 0644)
	if err != nil {
		log.Fatalln("cannot write into cache file: " + err.Error())
	}
}
