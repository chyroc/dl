package parse

import (
	"github.com/chyroc/dl/internal/resource"
)

type Parser interface {
	Kind() string
	Parse(uri string) (resource.Resource, error)
}
