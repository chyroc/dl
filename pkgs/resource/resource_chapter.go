package resource

import (
	"io"
)

type chapterResource struct {
	title       string
	chapterList []Resourcer
}

func NewChapter(title string, chapterList []Resourcer) Resourcer {
	return &chapterResource{
		title:       title,
		chapterList: chapterList,
	}
}

func (r *chapterResource) Title() string {
	return r.title
}

func (r *chapterResource) Reader() (int64, io.ReadCloser, error) {
	panic("un reachable")
}

func (r *chapterResource) Chapters() []Resourcer {
	return r.chapterList
}
