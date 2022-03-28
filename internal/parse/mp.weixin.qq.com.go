package parse

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chyroc/dl/internal/config"
	"github.com/chyroc/dl/internal/resource"
)

func NewMpWeixinQqCom() Parser {
	return &mpWeixinQqCom{}
}

type mpWeixinQqCom struct{}

func (r *mpWeixinQqCom) Kind() string {
	return "mp.weixin.qq.com"
}

func (r *mpWeixinQqCom) ExampleURLs() []string {
	return []string{
		"https://mp.weixin.qq.com/s/1Li7fFUu49XQoo6-dNwnmQ",
	}
}

func (r *mpWeixinQqCom) Parse(uri string) (resource.Resource, error) {
	text, err := config.ReqCli.New(http.MethodGet, uri).Text()
	if err != nil {
		return nil, err
	}

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(text))

	title := strings.TrimSpace(doc.Find("title").Text())
	if title == "" {
		title = strings.TrimSpace(doc.Find("meta[property='og:title']").AttrOr("content", ""))
	}
	if title == "" {
		title = strings.TrimSpace(doc.Find("meta[property='twitter:title']").AttrOr("content", ""))
	}

	vid := getMatchString(text, mpvidReg)
	mid := getMatchString(text, midReg)
	biz := getMatchString(text, bizReg)
	info, err := r.getVideoInfo(biz, mid, vid)
	if err != nil {
		return nil, err
	}

	spec := []*resource.Specification{}
	for _, v := range info.URLInfo {
		spec = append(spec, &resource.Specification{
			Size:       int64(v.Filesize),
			Definition: resource.MayConvertDefinition(v.VideoQualityWording),
			URL:        v.URL,
		})
	}

	return resource.NewURLChapter(title, []resource.Resource{
		resource.NewURLWithSpecification(info.Title+".mp4", spec),
	}), nil
}

var (
	mpvidReg = regexp.MustCompile(`data-mpvid="(.*?)"`)
	bizReg   = regexp.MustCompile(`var biz = "(.*?)"`)
	midReg   = regexp.MustCompile(`var mid = "(.*?)"`)
)

func (r *mpWeixinQqCom) getVideoInfo(biz, mid, vid string) (*getVideoInfoResp, error) {
	resp := new(getVideoInfoResp)
	url := fmt.Sprintf("https://mp.weixin.qq.com/mp/videoplayer?action=get_mp_video_play_url&__biz=%s&mid=%s&vid=%s&f=json", biz, mid, vid)
	err := config.ReqCli.New(http.MethodGet, url).Unmarshal(resp)
	return resp, err
}

type getVideoInfoResp struct {
	URLInfo []struct {
		URL                 string `json:"url"`
		FormatID            int    `json:"format_id"`
		DurationMs          int    `json:"duration_ms"`
		Filesize            int    `json:"filesize"`
		Width               int    `json:"width"`
		Height              int    `json:"height"`
		VideoQualityLevel   int    `json:"video_quality_level"`
		VideoQualityWording string `json:"video_quality_wording"`
	} `json:"url_info"`
	Title string `json:"title"`
}
