package download

import (
	"github.com/chyroc/dl/internal/helper"
)

func NewDownloadMp3(mp3 *helper.MP3, dir string, chapter bool) Downloader {
	return &downloadMp3{
		chapter: chapter,
		mp3:     mp3,
		dir:     dir,
	}
}

type downloadMp3 struct {
	dir     string
	chapter bool
	mp3     *helper.MP3
}

func (r *downloadMp3) Title() string {
	return helper.TrimFileExt(r.mp3.FileName)
}

func (r *downloadMp3) TargetFile() string {
	return r.mp3.FileName
}

func (r *downloadMp3) Download() error {
	file := r.TargetFile()
	if r.dir != "" {
		file = r.dir + "/" + file
	}
	return helper.DownloadMp32(r.mp3, r.chapter, file)
}

func (r *downloadMp3) MultiDownload() []Downloader {
	return nil
}
