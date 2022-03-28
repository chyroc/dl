package parse

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chyroc/dl/pkgs/config"
	"github.com/chyroc/dl/pkgs/resource"
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
		"https://mp.weixin.qq.com/s/xDNx_qGrrUoHK_HXBiKI5w",
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

	videoType := getMatchString(text, videoTypeReg)
	mid := getMatchString(text, midReg)
	biz := getMatchString(text, bizReg)

	switch videoType {
	case "1":
		videos, err := r.getType1Video(text)
		if err != nil {
			return nil, err
		}
		return resource.NewURLChapter(title, videos), nil
	case "2":
		videos, err := r.getType2Video(biz, mid, text)
		if err != nil {
			return nil, err
		}
		return resource.NewURLChapter(title, videos), nil
	default:
		return nil, fmt.Errorf("unknown video type: %s", videoType)
	}
}

func (r *mpWeixinQqCom) getType2Video(biz, mid, text string) (res []resource.Resource, err error) {
	matches := mpvidType2Reg.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		vid := match[1]
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

		res = append(res, resource.NewURLWithSpecification(info.Title+".mp4", spec))
	}
	return res, nil
}

func (r *mpWeixinQqCom) getType1Video(text string) (res []resource.Resource, err error) {
	matches := vidType1Reg.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if match[1] == "" {
			continue
		}
		vid := match[1]
		video, err := NewVQqCom().Parse(fmt.Sprintf("https://v.qq.com/txp/iframe/player.html?vid=%s", vid))
		if err != nil {
			return nil, err
		}

		res = append(res, video)
	}
	return res, nil
}

var (
	videoTypeReg  = regexp.MustCompile(`data-vidtype="(.*?)"`)
	vidType1Reg   = regexp.MustCompile(`vid=(.*?)"`)
	mpvidType2Reg = regexp.MustCompile(`data-mpvid="(.*?)"`)
	bizReg        = regexp.MustCompile(`var biz = "(.*?)"`)
	midReg        = regexp.MustCompile(`var mid = "(.*?)"`)
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
