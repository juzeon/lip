package source

import (
	"github.com/juzeon/lip/data"
	"github.com/juzeon/lip/util"
	"log"
	"net"
)

type ISource interface {
	Init() error
	Close()
	GetName() string
	GetDatabaseFileName() string
	DownloadDatabase() error
	LookUp(ip net.IP) (data.IPLookupResult, error)
	CheckUpdate() bool
	IsOnline() bool
}

var uninitializedSources = []ISource{
	&IP2Region{},
	&QQWry{},
	&IPApi{},
	&IPInfo{},
	&IPIP{},
}
var Sources []ISource

func InitDatabases() {
	for _, src := range uninitializedSources {
		src := src
		err := src.Init()
		if err != nil {
			log.Println("failed to initialize " + src.GetName() + ": " + err.Error())
			continue
		}
		Sources = append(Sources, src)
	}
}
func CloseDatabases() {
	for _, src := range Sources {
		src.Close()
	}
}
func DownloadDatabases(overwrite bool) {
	for _, src := range uninitializedSources {
		if src.IsOnline() {
			continue
		}
		dbPath := util.MustLipPath(src.GetDatabaseFileName())
		if !util.MustCheckExist(dbPath) || overwrite {
			log.Println("downloading database of " + src.GetName())
			err := src.DownloadDatabase()
			if err != nil {
				log.Println("failed to download " + src.GetName() + " into " + dbPath)
			} else {
				log.Println("downloaded " + src.GetName() + " successfully")
			}
		}
	}
}
