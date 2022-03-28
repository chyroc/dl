package common

import (
	"github.com/chyroc/dl/internal/resource"
)

type MusicRequest interface {
	Do() error
	Extract() ([]*resource.MP3, error)
}
