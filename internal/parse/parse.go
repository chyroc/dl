package parse

import (
	"github.com/chyroc/dl/internal/resource"
)

type Parser interface {
	Kind() string
	ExampleURLs() []string
	Parse(uri string) (resource.Resource, error)
}
