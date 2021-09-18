package parse

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/chyroc/dl/internal/config"
	"github.com/chyroc/dl/internal/download"
)

func NewWwwZhihuCom() Parser {
	return &wwwZhihuCom{}
}

type wwwZhihuCom struct{}

func (r *wwwZhihuCom) Kind() string {
	return "www.zhihu.com"
}

func (r *wwwZhihuCom) Parse(uri string) (download.Downloader, error) {
	text, err := config.ReqCli.New(http.MethodGet, uri).Text()
	if err != nil {
		return nil, err
	}

	title := getMatchString(text, wwwZhihuComTitleReg)
	id := getMatchString(text, wwwZhihuComIDReg)
	resp := new(wwwZhihuComGetResp)

	err = config.ReqCli.New(http.MethodGet, fmt.Sprintf("https://lens.zhihu.com/api/v4/videos/%s", id)).WithHeaders(prepareCommonHeader(uri, nil)).Unmarshal(resp)
	if err != nil {
		return nil, err
	}
	spec := []*download.Specification{}
	for _, v := range []wwwZhihuComDefi{resp.Playlist.Ld, resp.Playlist.Sd} {
		spec = append(spec, &download.Specification{
			Size:       v.Size,
			Definition: download.MayConvertDefinition(v.Format),
			URL:        v.PlayURL,
		})
	}
	return download.NewDownloadURL(title, title+".mp4", false, spec), nil
}

var (
	wwwZhihuComIDReg    = regexp.MustCompile(`src="https://video.zhihu.com/video/(.*?)\?`)
	wwwZhihuComTitleReg = regexp.MustCompile(`ZVideo-title">(.*?)<`)
)

type wwwZhihuComGetResp struct {
	Playlist struct {
		Ld wwwZhihuComDefi `json:"LD"`
		Sd wwwZhihuComDefi `json:"SD"`
	} `json:"playlist"`
	CoverURL   string `json:"cover_url"`
	PlaylistV2 struct {
		Sd wwwZhihuComDefi `json:"SD"`
	} `json:"playlist_v2"`
	Title interface{} `json:"title"`
}

type wwwZhihuComDefi struct {
	// Maxbitrate int     `json:"maxbitrate"`
	Format  string `json:"format"`
	PlayURL string `json:"play_url"`
	// Height     int     `json:"height"`
	// Channels   int     `json:"channels"`
	// Width      int     `json:"width"`
	// SampleRate int     `json:"sample_rate"`
	// Fps        int     `json:"fps"`
	// Duration   float64 `json:"duration"`
	// Bitrate    int     `json:"bitrate"`
	Size int64 `json:"size"`
}
