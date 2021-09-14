package parse

import (
	"github.com/chyroc/dl/internal/download"
)

type Parser interface {
	Parse(uri string) (download.Downloader, error)
}
