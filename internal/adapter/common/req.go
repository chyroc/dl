package common

import (
	"github.com/chyroc/dl/internal/helper"
)

type MusicRequest interface {
	Do() error
	Extract() ([]*helper.MP3, error)
}
