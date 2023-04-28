package storage

import (
	"github.com/juzeon/lip/util"
	"log"
	"os"
)

func init() {
	lipPath := util.MustLipPath()
	if !util.MustCheckExist(lipPath) {
		err := os.MkdirAll(lipPath, 0755)
		if err != nil {
			log.Fatalln("cannot mkdir in user home directory: " + err.Error())
		}
	}
}
