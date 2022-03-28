package netease

import (
	"fmt"
	"net/http"

	"github.com/chyroc/dl/pkgs/resource"
)

type DjradioParams struct {
	RadioId int `json:"radioId"`
	Limit   int `json:"limit"`
	Offset  int `json:"offset"`
}

type DjradioResponse struct {
	Code     int               `json:"code"`
	Msg      string            `json:"msg"`
	Programs []*DjradioProgram `json:"programs"`
}

type DjradioProgram struct {
	MainSong struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
	} `json:"mainSong"`
	Dj struct {
		Nickname  string `json:"nickname"`
		Signature string `json:"signature"`
		Brand     string `json:"brand"`
	} `json:"dj"`
	BlurCoverURL string `json:"blurCoverUrl"`
	Name         string `json:"name"`
	ID           int64  `json:"id"`
}

type DjradioRequest struct {
	Id       int
	Params   DjradioParams
	Response DjradioResponse
}

func NewDjradioRequest(id int) *DjradioRequest {
	return &DjradioRequest{Id: id, Params: DjradioParams{}}
}

func (s *DjradioRequest) Do() error {
	offset := 0
	limit := 30
	for {
		resp := new(DjradioResponse)
		if err := neteaseRequest(DjradioParams{
			RadioId: s.Id,
			Limit:   limit,
			Offset:  offset,
		}, DjradioAPI, resp); err != nil {
			return err
		} else if resp.Code != http.StatusOK {
			return fmt.Errorf(resp.Msg)
		}
		s.Response.Programs = append(s.Response.Programs, resp.Programs...)
		if len(resp.Programs) < limit {
			break
		}
		offset += limit
	}

	return nil
}

func (r *DjradioRequest) Extract() ([]*resource.MP3, error) {
	if len(r.Response.Programs) == 0 {
		return nil, nil
	}
	ids := make([]int, 0, len(r.Response.Programs))
	for _, i := range r.Response.Programs {
		ids = append(ids, i.MainSong.ID)
	}

	req := NewSongRequest(ids...)
	if err := req.Do(); err != nil {
		return nil, err
	}

	return ExtractMP3List(req.Response.Songs, r.Response.Programs[0].Dj.Brand)
}
