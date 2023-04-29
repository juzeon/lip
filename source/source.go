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
	LookUp(ip net.IP, duplicateIdentical bool) (data.IPLookupResult, error)
	CheckUpdate() bool
}

var uninitializedSources = []ISource{
	&IP2Region{},
	&QQWry{},
}
var Sources []ISource

func InitDatabases() {
	for _, ori := range uninitializedSources {
		ori := ori
		err := ori.Init()
		if err != nil {
			log.Println("failed to initialize " + ori.GetName())
			continue
		}
		Sources = append(Sources, ori)
	}
}
func DownloadDatabases(overwrite bool) {
	for _, ori := range uninitializedSources {
		dbPath := util.MustLipPath(ori.GetDatabaseFileName())
		if !util.MustCheckExist(dbPath) || overwrite {
			log.Println("downloading database of " + ori.GetName())
			err := ori.DownloadDatabase()
			if err != nil {
				log.Println("failed to download " + ori.GetName() + " into " + dbPath)
			}
			log.Println("downloaded " + ori.GetName() + " successfully")
		}
	}
}
