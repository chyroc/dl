package resource

import (
	"io"
	"net/http"
	"time"
)

type Resourcer interface {
	Title() string
	Reader() (int64, io.ReadCloser, error)
}

type Resourcer2 interface {
	Title() string
	Reader2() (func() int64, io.ReadCloser, error)
}

type ChapterResource interface {
	Resourcer
	Chapters() []Resourcer
}

type Mp3Resource interface {
	Resourcer
	MP3() *MP3
}

type MP3ChapterResource interface {
	Resourcer
	Chapters() []Mp3Resource
}

var downloadHttpClient = http.Client{
	Timeout: time.Second * 60,
}
