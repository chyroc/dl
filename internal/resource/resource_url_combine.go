package resource

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

type urlCombineResource struct {
	title    string
	urls     []string
	callback func(url string, data []byte) []byte
	sep      []byte
}

func NewURLCombineResource(title string, urls []string) Resource {
	return &urlCombineResource{
		title: title,
		urls:  urls,
	}
}

type URLCombineOption struct {
	Callback func(url string, data []byte) []byte
	Sep      []byte
}

func NewURLCombineWithOption(title string, urls []string, option *URLCombineOption) Resource {
	return &urlCombineResource{
		title:    title,
		urls:     urls,
		callback: option.Callback,
		sep:      option.Sep,
	}
}

func (r *urlCombineResource) Title() string {
	return r.title
}

func (r *urlCombineResource) Reader() (int64, io.ReadCloser, error) {
	data, err := r.fetchData()
	if err != nil {
		return 0, nil, err
	}
	return int64(len(data)), ioutil.NopCloser(bytes.NewReader(data)), nil
}

// break when return ErrCombineResourceEnd
func (r *urlCombineResource) fetchData() ([]byte, error) {
	var data [][]byte
	for _, url := range r.urls {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		resp, err := downloadHttpClient.Do(req)
		if err != nil {
			return nil, err
		}

		bs, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if r.callback != nil {
			bs = r.callback(url, bs)
		}

		data = append(data, bs)
	}

	return bytes.Join(data, r.sep), nil
}
