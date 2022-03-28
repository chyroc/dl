package netease

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chyroc/dl/internal/resource"
)

type SongUrlParams struct {
	Ids string `json:"ids"`
	Br  int    `json:"br"`
}

type SongUrlResponse struct {
	Code int       `json:"code"`
	Msg  string    `json:"msg"`
	Data []SongUrl `json:"data"`
}

type SongUrlRequest struct {
	Params   SongUrlParams
	Response SongUrlResponse
}

func NewSongUrlRequest(ids ...int) *SongUrlRequest {
	br := 128
	switch br {
	case 128, 192, 320:
		br *= 1000
		break
	default:
		br = 999 * 1000
	}
	enc, _ := json.Marshal(ids)
	return &SongUrlRequest{Params: SongUrlParams{Ids: string(enc), Br: br}}
}

func (s *SongUrlRequest) Do() error {
	if err := neteaseRequest(s.Params, SongUrlAPI, &s.Response); err != nil {
		return err
	} else if s.Response.Code != http.StatusOK {
		return fmt.Errorf(s.Response.Msg)
	}

	return nil
}

func (s *SongRequest) Do() error {
	if err := neteaseRequest(s.Params, SongAPI, &s.Response); err != nil {
		return err
	} else if s.Response.Code != http.StatusOK {
		return fmt.Errorf(s.Response.Msg)
	}

	return nil
}

func (s *SongRequest) Extract() ([]*resource.MP3, error) {
	return ExtractMP3List(s.Response.Songs, ".")
}
