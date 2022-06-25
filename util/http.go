package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// Downloader ...
type Downloader struct {
	URL     string
	Timeout time.Duration
}

func (d *Downloader) download() ([]byte, error) {
	request, err := http.NewRequest(http.MethodGet, d.URL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: d.Timeout}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("http_error [%d]", response.StatusCode)
	}

	return body, err
}

// SetTimeout ...
func (d *Downloader) SetTimeout(timeout time.Duration) *Downloader {
	d.Timeout = timeout
	return d
}

// DownloadText ...
func (d *Downloader) DownloadText() (string, error) {
	data, err := d.download()
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// DownloadToFile ...
func (d *Downloader) DownloadToFile(savePath, fileName string) error {
	data, err := d.download()
	if err != nil {
		return err
	}

	if len(savePath) == 0 {
		savePath = "./"
	}

	if len(fileName) == 0 {
		fileName = "util-downloader-default"
	}

	return ioutil.WriteFile(savePath+strings.Trim(fileName, "/"), data, os.FileMode(0666))
}

// Download ...
func (d *Downloader) Download() ([]byte, error) {
	return d.download()
}
