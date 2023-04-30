package httpclient

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
	"os"
	"strconv"
	"time"
)

var proxy string

func GetClient() *resty.Client {
	client := resty.New().
		SetTimeout(60 * time.Second).
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetHeaders(map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/112.0",
		})
	if proxy != "" {
		client.SetProxy(proxy)
	}
	return client
}
func SetProxy(url string) {
	proxy = url
}
func DownloadTo(url string, filePath string) error {
	_, err := GetClient().SetTimeout(5 * time.Second).R().Head(url)
	if err != nil {
		return err
	}
	resp, err := GetClient().R().Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return errors.New("respond with status code: " + strconv.Itoa(resp.StatusCode()))
	}
	err = os.WriteFile(filePath, resp.Body(), 0644)
	if err != nil {
		return err
	}
	return nil
}
func DownloadToWithMultipleURLs(urls []string, filePath string) error {
	var err error
	for _, url := range urls {
		log.Println("fetching data from: " + url)
		err = DownloadTo(url, filePath)
		if err != nil {
			log.Println("download failed: " + err.Error())
		} else {
			return nil
		}
	}
	return fmt.Errorf("download with multiple urls failed: %v", err)
}
