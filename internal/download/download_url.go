package download

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func NewDownloadURL(title string, pkgs []*Pkg) Downloader {
	return &downloadURL{
		title: title,
		pkgs:  pkgs,
	}
}

type downloadURL struct {
	title string
	pkgs  []*Pkg
}

func (r *downloadURL) Title() string {
	return r.title
}

func (r *downloadURL) Download() error {
	pkg := PkgList(r.pkgs).GetMax()

	tempFile, err := downloadURL2(pkg.URL, nil)
	if err != nil {
		return err
	}
	return os.Rename(tempFile, r.Title()+".mp4")
}

type Pkg struct {
	Size       int64      `json:"size"`
	Definition Definition `json:"definition"`
	URL        string     `json:"url"`
	Type       string     `json:"type"`
}

type PkgList []*Pkg

func (r PkgList) GetMax() *Pkg {
	var size int64
	var pkg *Pkg
	for _, v := range r {
		if v.Size > size {
			size = v.Size
			pkg = v
		}
	}
	return pkg
}

// https://www.image-engineering.de/library/technotes/991-separating-sd-hd-full-hd-4k-and-8k
type Definition string

const (
	DefinitionSD     Definition = "sd"
	DefinitionHD     Definition = "hd"
	DefinitionFullHD Definition = "full-hd"
	DefinitionUHD    Definition = "uhd"
	Definition4K     Definition = "4k"
	Definition8K     Definition = "8k"
)

func downloadURL2(uri string, headers map[string]string) (string, error) {
	f, err := ioutil.TempFile("", "dl-*")
	if err != nil {
		return "", err
	}
	defer f.Close()

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return "", err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := downloadHttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	length := resp.Header.Get("Content-Length")
	_ = length

	for {
		buf := make([]byte, 1024)
		n, err := resp.Body.Read(buf)
		if n > 0 {
			if _, err := f.Write(buf[:n]); err != nil {
				return "", err
			}
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return "", err
		}
	}
	return f.Name(), nil
}

var downloadHttpClient = http.Client{}
