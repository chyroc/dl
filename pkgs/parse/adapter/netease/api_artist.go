package netease

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/chyroc/dl/pkgs/helper"
	"github.com/chyroc/dl/pkgs/resource"
)

type ArtistParams struct{}

type ArtistResponse struct {
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
	Artist   Artist `json:"artist"`
	HotSongs []Song `json:"hotSongs"`
}

type ArtistRequest struct {
	Id       int
	Params   ArtistParams
	Response ArtistResponse
}

func NewArtistRequest(id int) *ArtistRequest {
	return &ArtistRequest{Id: id, Params: ArtistParams{}}
}

func (r *ArtistRequest) Do() error {
	if err := neteaseRequest(r.Params, ArtistAPI+fmt.Sprintf("/%d", r.Id), &r.Response); err != nil {
		return err
	} else if r.Response.Code != http.StatusOK {
		return fmt.Errorf(r.Response.Msg)
	}

	return nil
}

func (r *ArtistRequest) Extract() ([]*resource.MP3, error) {
	ids := make([]int, 0, len(r.Response.HotSongs))
	for _, i := range r.Response.HotSongs {
		ids = append(ids, i.Id)
	}

	req := NewSongRequest(ids...)
	if err := req.Do(); err != nil {
		return nil, err
	}

	savePath := filepath.Join(".", helper.TrimInvalidFilePathChars(r.Response.Artist.Name))
	return ExtractMP3List(req.Response.Songs, savePath)
}
