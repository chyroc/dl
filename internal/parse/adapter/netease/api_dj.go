package netease

import (
	"fmt"
	"net/http"

	"github.com/chyroc/dl/internal/resource"
)

type DJParams struct {
	ID int `json:"id"`
}

type DJResponse struct {
	Code    int             `json:"code"`
	Msg     string          `json:"msg"`
	Program *DjradioProgram `json:"program"`
}

type DJRequest struct {
	Params   DJParams
	Response DJResponse
}

func NewDJRequest(id int) *DJRequest {
	return &DJRequest{Params: DJParams{ID: id}}
}

// result['program']['mainDJ']['id'],result['program']['mainDJ']['name']
func (s *DJRequest) Do() error {
	if err := neteaseRequest(s.Params, DJAPI, &s.Response); err != nil {
		return err
	} else if s.Response.Code != http.StatusOK {
		return fmt.Errorf(s.Response.Msg)
	}

	return nil
}

func (r *DJRequest) Extract() ([]*resource.MP3, error) {
	req := NewSongRequest(r.Response.Program.MainSong.ID)
	if err := req.Do(); err != nil {
		return nil, err
	}

	return ExtractMP3List(req.Response.Songs, r.Response.Program.Dj.Brand)
}
