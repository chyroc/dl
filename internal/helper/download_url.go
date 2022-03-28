package helper

import (
	"io"
	"io/ioutil"
	"net/http"
)

var downloadHttpClient = http.Client{}

func Download(url string, chapter bool) (string, error) {
	f, err := ioutil.TempFile("", "dl-*")
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := downloadHttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	reader := NewProgressReaderClose(resp.ContentLength, resp.Body, chapter)
	defer reader.Close()

	if _, err := io.Copy(f, reader); err != nil {
		return "", err
	}

	return f.Name(), nil
}

func Download2(url, target string, chapter bool) error {
	file, err := Download(url, chapter)
	if err != nil {
		return err
	}
	return Rename(file, target)
}
