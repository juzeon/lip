package httpclient

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"os"
	"strconv"
	"time"
)

var client *resty.Client

func init() {
	client = resty.New().SetTimeout(120 * time.Second)
}
func SetProxy(url string) {
	client = client.SetProxy(url)
}
func DownloadTo(url string, filePath string) error {
	resp, err := client.R().Get(url)
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
