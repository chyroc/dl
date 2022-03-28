package resource

import (
	"io"
	"net/http"
)

type Resource interface {
	Title() string
	Reader() (int64, io.ReadCloser, error)
}

type ChapterResource interface {
	Resource
	ChapterCount() int
	Chapters() []Resource
}

type Mp3Resource interface {
	Resource
	MP3() *MP3
}

var downloadHttpClient = http.Client{}
