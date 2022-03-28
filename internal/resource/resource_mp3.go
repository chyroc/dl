package resource

import (
	"io"
	"net/http"
)

type mp3Resource struct {
	mp3 *MP3
}

func NewMp3(mp3 *MP3) Resource {
	return &mp3Resource{
		mp3: mp3,
	}
}

func NewMp32(mp3 *MP3) Mp3Resource {
	return &mp3Resource{
		mp3: mp3,
	}
}

func (r *mp3Resource) Title() string {
	return r.mp3.FileName
}

func (r *mp3Resource) Reader() (int64, io.ReadCloser, error) {
	url := r.mp3.DownloadUrl

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, nil, err
	}

	resp, err := downloadHttpClient.Do(req)
	if err != nil {
		return 0, nil, err
	}

	return resp.ContentLength, resp.Body, nil
}

func (r *mp3Resource) MP3() *MP3 {
	return r.mp3
}
