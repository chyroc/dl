package download

import (
	"github.com/chyroc/dl/internal/helper"
)

func NewDownloadMp3List(dir, album string, mp3List []*helper.MP3) Downloader {
	return &downloadMp3List{
		dir:     dir,
		album:   album,
		Mp3List: mp3List,
	}
}

type downloadMp3List struct {
	dir     string
	album   string
	Mp3List []*helper.MP3
}

func (r *downloadMp3List) Title() string {
	return helper.TrimFileExt(r.album)
}

func (r *downloadMp3List) TargetFile() string {
	return ""
}

func (r *downloadMp3List) Download() error {
	return nil
}

func (r *downloadMp3List) MultiDownload() (res []Downloader) {
	for _, mp3 := range r.Mp3List {
		res = append(res, NewDownloadMp3(mp3, r.dir, true))
	}
	return res
}
