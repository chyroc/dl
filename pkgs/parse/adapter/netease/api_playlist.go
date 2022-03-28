package netease

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/chyroc/dl/pkgs/helper"
	"github.com/chyroc/dl/pkgs/resource"
)

type PlaylistParams struct {
	Id int `json:"id"`
}

type PlaylistResponse struct {
	Code     int      `json:"code"`
	Msg      string   `json:"msg"`
	Playlist Playlist `json:"playlist"`
}

type PlaylistRequest struct {
	Params   PlaylistParams
	Response PlaylistResponse
}

func NewPlaylistRequest(id int) *PlaylistRequest {
	return &PlaylistRequest{Params: PlaylistParams{Id: id}}
}

func (s *PlaylistRequest) Do() error {
	if err := neteaseRequest(s.Params, PlaylistAPI, &s.Response); err != nil {
		return err
	} else if s.Response.Code != http.StatusOK {
		return fmt.Errorf(s.Response.Msg)
	}

	return nil
}

func (p *PlaylistRequest) Extract() ([]*resource.MP3, error) {
	ids := make([]int, 0, len(p.Response.Playlist.TrackIds))
	for _, i := range p.Response.Playlist.TrackIds {
		ids = append(ids, i.Id)
	}

	req := NewSongRequest(ids...)
	if err := req.Do(); err != nil {
		return nil, err
	}

	savePath := filepath.Join(".", helper.TrimInvalidFilePathChars(p.Response.Playlist.Name))
	return ExtractMP3List(req.Response.Songs, savePath)
}
