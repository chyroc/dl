package parse

import (
	"fmt"

	"github.com/chyroc/dl/internal/adapter/netease"
	"github.com/chyroc/dl/internal/download"
)

func NewMusic163Com() Parser {
	return &music163Com{}
}

type music163Com struct{}

func (r *music163Com) Kind() string {
	return "music.163.com"
}

func (r *music163Com) Parse(uri string) (download.Downloader, error) {
	req, err := netease.Parse(uri)
	if err != nil {
		return nil, err
	}

	if err = req.Do(); err != nil {
		return nil, err
	}

	mp3List, err := req.Extract()
	if err != nil {
		return nil, err
	}

	if len(mp3List) == 0 {
		return nil, fmt.Errorf("find no mp3")
	}
	if len(mp3List) == 1 {
		return download.NewDownloadMp3(mp3List[0], "", false), nil
	}
	return download.NewDownloadMp3List(mp3List[0].Tag.Album, mp3List[0].Tag.Album, mp3List), nil
}
