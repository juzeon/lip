package storage

import (
	"github.com/juzeon/lip/httpclient"
	"github.com/juzeon/lip/origin"
	"github.com/juzeon/lip/util"
	"log"
)

func DownloadDatabases(overwrite bool) {
	for _, ori := range origin.Origins {
		dbPath := util.MustLipPath(ori.DatabaseFileName)
		if !util.MustCheckExist(dbPath) || overwrite {
			log.Println("downloading " + ori.Name)
			var err error
			for _, url := range ori.DatabaseURLs {
				log.Println("try url: " + url)
				err = httpclient.DownloadTo(url, dbPath)
				if err != nil {
					log.Println("download failed: " + err.Error())
				} else {
					break
				}
			}
			if err != nil {
				log.Fatalln("failed to download " + ori.Name + " into " + dbPath)
			}
		}
	}
}
