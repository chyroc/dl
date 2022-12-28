package parse

import (
	"strings"

	"github.com/chyroc/dl/pkgs/resource"
)

func NewM3u8() Parser {
	return &m3u8{}
}

type m3u8 struct{}

func (r *m3u8) Kind() string {
	return "filetype.m3u8"
}

func (r *m3u8) ExampleURLs() []string {
	return []string{""}
}

func (r *m3u8) Parse(uri string) (resource.Resourcer, error) {
	x := strings.Split(uri, ".")
	base := x[len(x)-1]
	return resource.NewM3U8(100, base, uri)
}
