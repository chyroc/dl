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
	return 0, newReaderByMultiURLs(r.urls, r.callback, r.sep), nil
}

func newReaderByMultiURLs(urls []string, callback func(url string, data []byte) []byte, sep []byte) io.ReadCloser {
	return &readerByMultiURLs{urls: urls, callback: callback, sep: sep}
}

type readerByMultiURLs struct {
	urls     []string
	sep      []byte
	callback func(url string, data []byte) []byte

	fetched bool
	data    []byte
	dataIdx int
}

func (r *readerByMultiURLs) Read(p []byte) (n int, err error) {
	if !r.fetched {
		textList := [][]byte{}

		for _, contentURL := range r.urls {
			resp, err := http.Get(contentURL)
			if err != nil {
				return 0, err
			}
			bs, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return 0, err
			}
			if r.callback != nil {
				textList = append(textList, r.callback(contentURL, bs))
			} else {
				textList = append(textList, bs)
			}
		}

		if len(r.sep) > 0 {
			r.data = bytes.Join(textList, r.sep)
		}
		r.fetched = true
	}

	if r.dataIdx >= len(r.data) {
		return 0, io.EOF
	}
	r.dataIdx = copy(p, r.data[r.dataIdx:])
	return r.dataIdx, nil
}

func (r *readerByMultiURLs) Close() error {
	return nil
}
