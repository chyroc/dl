package netease

import (
	"encoding/json"
	"net/http"

	"github.com/chyroc/dl/internal/config"
)

const (
	WeAPI       = "https://music.163.com/weapi"
	SongUrlAPI  = "https://music.163.com/weapi/song/enhance/player/url"
	SongAPI     = "https://music.163.com/weapi/v3/song/detail"
	ArtistAPI   = "https://music.163.com/weapi/v1/artist"
	AlbumAPI    = "https://music.163.com/weapi/v1/album"
	PlaylistAPI = "https://music.163.com/weapi/v3/playlist/detail"
	DjradioAPI  = "https://music.163.com/weapi/dj/program/byradio"
	DJAPI       = "https://music.163.com/weapi/dj/program/detail"
)

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

	text, err := config.ReqCli.New(http.MethodPost, url).WithFormURLEncoded(body).WithHeaders(header).Text()
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(text), resp); err != nil {
		return err
	}
	return nil
}
