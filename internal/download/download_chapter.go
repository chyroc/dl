package download

import (
	"fmt"
)

type Chapter struct {
	Tid   string
	Pid   string
	P     int
	Title string
	URL   string
}

func NewDownloadChapter(dir, title string, ext string, chapters []*Chapter) Downloader {
	return &downloadChapter{
		dir:      dir,
		title:    title,
		ext:      ext,
		chapters: chapters,
	}
}

type downloadChapter struct {
	dir      string
	title    string
	ext      string
	chapters []*Chapter
}

func (r *downloadChapter) Title() string {
	return r.title
}

func (r *downloadChapter) TargetFile() string {
	return ""
}

func (r *downloadChapter) Download() error {
	return nil
}

func (r *downloadChapter) MultiDownload() (res []Downloader) {
	for _, down := range r.chapters {
		res = append(res, NewDownloadURL(fmt.Sprintf("[%s] %s", down.Pid, down.Title), r.dir+"/"+down.Title+r.ext, true, []*Specification{{URL: down.URL}}))
	}
	return res
}
