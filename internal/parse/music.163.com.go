package parse

import (
	"fmt"

	"github.com/chyroc/dl/internal/adapter/netease"
	"github.com/chyroc/dl/internal/resource"
)

func NewMusic163Com() Parser {
	return &music163Com{}
}

type music163Com struct{}

func (r *music163Com) Kind() string {
	return "music.163.com"
}

func (r *music163Com) ExampleURLs() []string {
	return []string{
		"https://music.163.com/#/song?id=1843572582",
		"https://music.163.com/#/album?id=132874562",
		"https://music.163.com/#/playlist?id=156934569",
	}
}

func (r *music163Com) Parse(uri string) (resource.Resource, error) {
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
		return resource.NewMp3(mp3List[0]), nil
	}
	panic("")
	// return download.NewDownloadMp3List(mp3List[0].SavePath, mp3List[0].SavePath, mp3List), nil
}
