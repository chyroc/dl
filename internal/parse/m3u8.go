package parse

import (
	"strings"

	"github.com/chyroc/dl/internal/resource"
)

func NewM3u8() Parser {
	return &m3u8{}
}

type m3u8 struct{}

func (r *m3u8) Kind() string {
	return "filetype.m3u8"
}

func (r *m3u8) Parse(uri string) (resource.Resource, error) {
	x := strings.Split(uri, ".")
	base := x[len(x)-1]
	return resource.NewURL(base, uri), nil
}
