package tencent

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/chyroc/dl/internal/config"
	"github.com/chyroc/dl/internal/helper"
)

const (
	SongAPI     = "https://c.y.qq.com/v8/fcg-bin/fcg_play_single_song.fcg"
	ArtistAPI   = "https://c.y.qq.com/v8/fcg-bin/fcg_v8_singer_track_cp.fcg"
	AlbumAPI    = "https://c.y.qq.com/v8/fcg-bin/fcg_v8_album_detail_cp.fcg"
	PlaylistAPI = "https://c.y.qq.com/v8/fcg-bin/fcg_v8_playlist_cp.fcg"
)

type SongResponse struct {
	Code int    `json:"code"`
	Data []Song `json:"data"`
}

type SongRequest struct {
	Params   map[string]string
	Response SongResponse
}

type SingerResponse struct {
	Code int `json:"code"`
	Data struct {
		List []struct {
			MusicData Song `json:"musicData"`
		} `json:"list"`
		SingerId   string `json:"singer_id"`
		SingerMid  string `json:"singer_mid"`
		SingerName string `json:"singer_name"`
		Total      int    `json:"total"`
	} `json:"data"`
}

type SingerRequest struct {
	Params   map[string]string
	Response SingerResponse
}

type AlbumResponse struct {
	Code int
	Data struct {
		GetAlbumInfo GetAlbumInfo `json:"getAlbumInfo"`
		GetSongInfo  []Song       `json:"getSongInfo"`
	} `json:"data"`
}

type AlbumRequest struct {
	Params   map[string]string
	Response AlbumResponse
}

type PlaylistResponse struct {
	Code int `json:"code"`
	Data struct {
		CDList []CD `json:"cdlist"`
	} `json:"data"`
}

type PlaylistRequest struct {
	Params   map[string]string
	Response PlaylistResponse
}

func NewSongRequest(mid string) *SongRequest {
	query := map[string]string{
		"songmid":  mid,
		"platform": "yqq",
		"format":   "json",
	}
	return &SongRequest{Params: query}
}

func NewSingerRequest(mid string) *SingerRequest {
	query := map[string]string{
		"singermid": mid,
		"begin":     "0",
		"num":       "50",
		"order":     "listen",
		"newsong":   "1",
		"platform":  "yqq",
	}
	return &SingerRequest{Params: query}
}

func NewAlbumRequest(mid string) *AlbumRequest {
	query := map[string]string{
		"albummid": mid,
		"newsong":  "1",
		"platform": "yqq",
		"format":   "json",
	}
	return &AlbumRequest{Params: query}
}

func NewPlaylistRequest(id string) *PlaylistRequest {
	query := map[string]string{
		"id":       id,
		"newsong":  "1",
		"platform": "yqq",
		"format":   "json",
	}
	return &PlaylistRequest{Params: query}
}

func (s *SongRequest) Do() error {
	err := tencentRequest(s.Params, SongAPI, &s.Response)
	if err != nil {
		return err
	} else if s.Response.Code != 0 {
		return fmt.Errorf("%d", s.Response.Code)
	}

	return nil
}

func (s *SongRequest) Extract() ([]*helper.MP3, error) {
	return ExtractMP3List(s.Response.Data, ".")
}

func (s *SingerRequest) Do() error {
	err := tencentRequest(s.Params, ArtistAPI, &s.Response)
	if err != nil {
		return err
	} else if s.Response.Code != 0 {
		return fmt.Errorf("%d", s.Response.Code)
	}
	return nil
}

func (a *SingerRequest) Extract() ([]*helper.MP3, error) {
	savePath := filepath.Join(".", helper.TrimInvalidFilePathChars(a.Response.Data.SingerName))
	var songs []Song
	for _, i := range a.Response.Data.List {
		songs = append(songs, i.MusicData)
	}
	return ExtractMP3List(songs, savePath)
}

func (s *AlbumRequest) Do() error {
	err := tencentRequest(s.Params, AlbumAPI, &s.Response)
	if err != nil {
		return err
	} else if s.Response.Code != 0 {
		return fmt.Errorf("%d", s.Response.Code)
	}

	return nil
}

func (a *AlbumRequest) Extract() ([]*helper.MP3, error) {
	savePath := filepath.Join(".", helper.TrimInvalidFilePathChars(a.Response.Data.GetAlbumInfo.FAlbumName))
	return ExtractMP3List(a.Response.Data.GetSongInfo, savePath)
}

func (s *PlaylistRequest) Do() error {
	err := tencentRequest(s.Params, PlaylistAPI, &s.Response)
	if err != nil {
		return err
	} else if s.Response.Code != 0 {
		return fmt.Errorf("%d", s.Response.Code)
	}

	return nil
}

func (p *PlaylistRequest) Extract() ([]*helper.MP3, error) {
	var res []*helper.MP3
	for _, i := range p.Response.Data.CDList {
		savePath := filepath.Join(".", helper.TrimInvalidFilePathChars(i.DissName))
		mp3List, err := ExtractMP3List(i.SongList, savePath)
		if err != nil {
			continue
		}
		res = append(res, mp3List...)
	}

	return res, nil
}

func prepareCommonHeader(uri string, s map[string]string) map[string]string {
	// uriParsed, _ := url.Parse(uri)
	res := map[string]string{
		// "Host":       uriParsed.Host,
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36",
		"Accept":     "*/*",
		"Origin":     uri,
		"Referer":    uri,
	}
	for k, v := range s {
		res[k] = v
	}
	return res
}

func tencentRequest(query map[string]string, url string, resp interface{}) error {
	header := prepareCommonHeader("https://c.y.qq.com", nil)
	return config.ReqCli.New(http.MethodGet, url).WithQuerys(query).WithHeaders(header).Unmarshal(resp)
}
