package parse

import (
	"github.com/chyroc/dl/pkgs/resource"
)

type Parser interface {
	Kind() string
	ExampleURLs() []string
	Parse(uri string) (resource.Resourcer, error)
}
