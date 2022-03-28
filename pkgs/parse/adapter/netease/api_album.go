package netease

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/chyroc/dl/pkgs/helper"
	"github.com/chyroc/dl/pkgs/resource"
)

type AlbumParams struct{}

type AlbumResponse struct {
	SongResponse
	Album Album `json:"album"`
}

type AlbumRequest struct {
	Id       int
	Params   AlbumParams
	Response AlbumResponse
}

func NewAlbumRequest(id int) *AlbumRequest {
	return &AlbumRequest{Id: id, Params: AlbumParams{}}
}

func (s *AlbumRequest) Do() error {
	if err := neteaseRequest(s.Params, AlbumAPI+fmt.Sprintf("/%d", s.Id), &s.Response); err != nil {
		return err
	} else if s.Response.Code != http.StatusOK {
		return fmt.Errorf(s.Response.Msg)
	}

	return nil
}

func (a *AlbumRequest) Extract() ([]*resource.MP3, error) {
	savePath := filepath.Join(".", helper.TrimInvalidFilePathChars(a.Response.Album.Name))
	for i := range a.Response.Songs {
		a.Response.Songs[i].PublishTime = a.Response.Album.PublishTime
	}
	return ExtractMP3List(a.Response.Songs, savePath)
}
