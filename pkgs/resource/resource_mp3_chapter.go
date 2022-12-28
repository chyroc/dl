package resource

import (
	"io"
)

type mp3Chapter struct {
	title       string
	chapterList []Mp3Resource
}

func NewMP3Chapter(title string, chapterList []Mp3Resource) Resourcer {
	return &mp3Chapter{
		title:       title,
		chapterList: chapterList,
	}
}

func (r *mp3Chapter) Title() string {
	return r.title
}

func (r *mp3Chapter) Reader() (int64, io.ReadCloser, error) {
	panic("un reachable")
}

func (r *mp3Chapter) Chapters() []Mp3Resource {
	return r.chapterList
}
