package netease

import (
	"encoding/json"
)

type SongParams struct {
	C string `json:"c"`
}

type SongResponse struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Songs []Song `json:"songs"`
}

type SongRequest struct {
	Params   SongParams
	Response SongResponse
}

func NewSongRequest(ids ...int) *SongRequest {
	c := make([]map[string]int, 0, len(ids))
	for _, id := range ids {
		c = append(c, map[string]int{"id": id})
	}

	enc, _ := json.Marshal(c)
	return &SongRequest{Params: SongParams{C: string(enc)}}
}
