package netease

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/chyroc/dl/internal/config"
	"github.com/chyroc/dl/internal/helper"
)

const (
	WeAPI       = "https://music.163.com/weapi"
	SongUrlAPI  = "https://music.163.com/weapi/song/enhance/player/url"
	SongAPI     = "https://music.163.com/weapi/v3/song/detail"
	ArtistAPI   = "https://music.163.com/weapi/v1/artist"
	AlbumAPI    = "https://music.163.com/weapi/v1/album"
	PlaylistAPI = "https://music.163.com/weapi/v3/playlist/detail"
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

func NewSongRequest(ids ...int) *SongRequest {
	c := make([]map[string]int, 0, len(ids))
	for _, id := range ids {
		c = append(c, map[string]int{"id": id})
	}

	enc, _ := json.Marshal(c)
	return &SongRequest{Params: SongParams{C: string(enc)}}
}

func NewArtistRequest(id int) *ArtistRequest {
	return &ArtistRequest{Id: id, Params: ArtistParams{}}
}

func NewAlbumRequest(id int) *AlbumRequest {
	return &AlbumRequest{Id: id, Params: AlbumParams{}}
}

func NewPlaylistRequest(id int) *PlaylistRequest {
	return &PlaylistRequest{Params: PlaylistParams{Id: id}}
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

func (s *SongRequest) Extract() ([]*helper.MP3, error) {
	return ExtractMP3List(s.Response.Songs, ".")
}

func (s *ArtistRequest) Do() error {
	if err := neteaseRequest(s.Params, ArtistAPI+fmt.Sprintf("/%d", s.Id), &s.Response); err != nil {
		return err
	} else if s.Response.Code != http.StatusOK {
		return fmt.Errorf(s.Response.Msg)
	}

	return nil
}

func (a *ArtistRequest) Extract() ([]*helper.MP3, error) {
	ids := make([]int, 0, len(a.Response.HotSongs))
	for _, i := range a.Response.HotSongs {
		ids = append(ids, i.Id)
	}

	req := NewSongRequest(ids...)
	if err := req.Do(); err != nil {
		return nil, err
	}

	savePath := filepath.Join(".", helper.TrimInvalidFilePathChars(a.Response.Artist.Name))
	return ExtractMP3List(req.Response.Songs, savePath)
}

func (s *AlbumRequest) Do() error {
	if err := neteaseRequest(s.Params, AlbumAPI+fmt.Sprintf("/%d", s.Id), &s.Response); err != nil {
		return err
	} else if s.Response.Code != http.StatusOK {
		return fmt.Errorf(s.Response.Msg)
	}

	return nil
}

func (a *AlbumRequest) Extract() ([]*helper.MP3, error) {
	savePath := filepath.Join(".", helper.TrimInvalidFilePathChars(a.Response.Album.Name))
	for i := range a.Response.Songs {
		a.Response.Songs[i].PublishTime = a.Response.Album.PublishTime
	}
	return ExtractMP3List(a.Response.Songs, savePath)
}

func (s *PlaylistRequest) Do() error {
	if err := neteaseRequest(s.Params, PlaylistAPI, &s.Response); err != nil {
		return err
	} else if s.Response.Code != http.StatusOK {
		return fmt.Errorf(s.Response.Msg)
	}

	return nil
}

func (p *PlaylistRequest) Extract() ([]*helper.MP3, error) {
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

func neteaseRequest(param interface{}, url string, resp interface{}) error {
	enc, _ := json.Marshal(param)
	params, encSecKey, err := Encrypt(enc)
	if err != nil {
		return err
	}

	body := map[string]string{
		"params":    params,
		"encSecKey": encSecKey,
	}
	header := prepareCommonHeader(WeAPI, nil)

	if err := config.ReqCli.New(http.MethodPost, url).WithFormURLEncoded(body).WithHeaders(header).Unmarshal(resp); err != nil {
		return err
	}
	return nil
}
