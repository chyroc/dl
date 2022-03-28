package parse

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/chyroc/dl/pkgs/config"
	"github.com/chyroc/dl/pkgs/resource"
)

func NewWwwZhihuCom() Parser {
	return &wwwZhihuCom{}
}

type wwwZhihuCom struct{}

func (r *wwwZhihuCom) Kind() string {
	return "www.zhihu.com"
}

func (r *wwwZhihuCom) ExampleURLs() []string {
	return []string{""}
}

func (r *wwwZhihuCom) Parse(uri string) (resource.Resource, error) {
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
	spec := []*resource.Specification{}
	for _, v := range []wwwZhihuComDefi{resp.Playlist.Ld, resp.Playlist.Sd} {
		spec = append(spec, &resource.Specification{
			Size:       v.Size,
			Definition: resource.MayConvertDefinition(v.Format),
			URL:        v.PlayURL,
		})
	}
	return resource.NewURLWithSpecification(title+".mp4", spec), nil
}

var (
	wwwZhihuComIDReg    = regexp.MustCompile(`"videoId":"(.*?)"`)
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
