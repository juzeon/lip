package util

import (
	"errors"
	"log"
	"os"
	"path"
)

func MustLipPath(pathSeg ...string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("cannot get the path of user home directory: " + err.Error())
	}
	var arr []string
	arr = append(arr, home, ".lip")
	arr = append(arr, pathSeg...)
	return path.Join(arr...)
}
func MustCheckExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) {
		return false
	} else if err != nil {
		log.Fatalln("cannot access to lip's path: " + err.Error())
	}
	return true
}
