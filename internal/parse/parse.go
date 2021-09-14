package parse

import (
	"github.com/chyroc/dl/internal/download"
)

type Parser interface {
	Kind() string
	Parse(uri string) (download.Downloader, error)
}
