package parse

import (
	"strings"

	"github.com/chyroc/dl/internal/download"
)

func NewM3u8() Parser {
	return &m3u8{}
}

type m3u8 struct{}

func (r *m3u8) Kind() string {
	return "filetype.m3u8"
}

func (r *m3u8) Parse(uri string) (download.Downloader, error) {
	x := strings.Split(uri, ".")
	base := x[len(x)-1]
	return download.NewDownloadM3U8(base, base, uri), nil
}
