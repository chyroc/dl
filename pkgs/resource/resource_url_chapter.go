package resource

import (
	"io"
)

type urlChapter struct {
	title       string
	chapterList []Resource
}

type Chapter struct {
	ID    string
	Title string
	URL   string
}

func NewURLChapter(title string, chapterList []Resource) Resource {
	return &urlChapter{
		title:       title,
		chapterList: chapterList,
	}
}

func (r *urlChapter) Title() string {
	return r.title
}

func (r *urlChapter) Reader() (int64, io.ReadCloser, error) {
	panic("un reachable")
}

func (r *urlChapter) Chapters() []Resource {
	return r.chapterList
}
