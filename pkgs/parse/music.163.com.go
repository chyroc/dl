package parse

import (
	"fmt"

	"github.com/chyroc/dl/pkgs/parse/adapter/netease"
	"github.com/chyroc/dl/pkgs/resource"
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
		"https://music.163.com/#/discover/toplist?id=1978921795",
		"https://music.163.com/#/artist?id=905705",
		"https://music.163.com/#/djradio?id=970764541",
		"https://music.163.com/#/program?id=2499619396",
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

	chapters := []resource.Mp3Resource{}
	for _, v := range mp3List {
		chapters = append(chapters, resource.NewMp32(v))
	}

	switch req := req.(type) {
	case *netease.SongRequest:
		panic("unreachable song multi data")
	case *netease.ArtistRequest:
		return resource.NewMP3Chapter(req.Response.Artist.Name, chapters), nil
	case *netease.AlbumRequest:
		return resource.NewMP3Chapter(req.Response.Album.Name, chapters), nil
	case *netease.PlaylistRequest:
		return resource.NewMP3Chapter(req.Response.Playlist.Name, chapters), nil
	case *netease.DjradioRequest:
		return resource.NewMP3Chapter(req.Response.Programs[0].Dj.Brand, chapters), nil
	default:
		panic(fmt.Sprintf("unreachable type: %T", req))
	}
}
