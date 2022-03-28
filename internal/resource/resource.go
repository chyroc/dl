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
	Chapters() []Resource
}

type Mp3Resource interface {
	Resource
	MP3() *MP3
}

type MP3ChapterResource interface {
	Resource
	Chapters() []Mp3Resource
}

var downloadHttpClient = http.Client{}
