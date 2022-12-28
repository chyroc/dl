package parse

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/chyroc/dl/pkgs/config"
	"github.com/chyroc/dl/pkgs/resource"
)

func NewVCctvCom() Parser {
	return &vCctvCom{}
}

type vCctvCom struct{}

func (r *vCctvCom) Kind() string {
	return "v.cctv.com"
}

func (r *vCctvCom) ExampleURLs() []string {
	return []string{"https://v.cctv.com/2021/09/17/VIDERZvtKr1arx2zGkZprwqR210917.shtml"}
}

func (r *vCctvCom) Parse(uri string) (resource.Resourcer, error) {
	guid, err := r.getVideoID(uri)
	if err != nil {
		return nil, err
	}
	meta, err := r.getMeta(uri, guid)
	if err != nil {
		return nil, err
	}

	return resource.NewURL(meta.Title+".mp4", meta.HlsURL), nil
}

func (r *vCctvCom) getVideoID(uri string) (string, error) {
	text, err := config.ReqCli.New(http.MethodGet, uri).Text()
	if err != nil {
		return "", err
	}
	return getMatchString(text, vCctvComGuidReg), nil
}

func (r *vCctvCom) getMeta(originURL, guid string) (*vCctvComGetMetaResp, error) {
	uri := fmt.Sprintf("https://vdn.apps.cntv.cn/api/getHttpVideoInfo.do?pid=%s", guid)
	header := prepareCommonHeader(originURL, nil)
	resp := new(vCctvComGetMetaResp)
	err := config.ReqCli.New(http.MethodGet, uri).WithHeaders(header).Unmarshal(resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type vCctvComGetMetaResp struct {
	Title string `json:"title"`
	Video struct {
		TotalLength string `json:"totalLength"`
		Chapters    []struct {
			Duration string `json:"duration"`
			Image    string `json:"image"`
			URL      string `json:"url"`
		} `json:"chapters"`
		Chapters2 []struct {
			Duration string `json:"duration"`
			Image    string `json:"image"`
			URL      string `json:"url"`
		} `json:"chapters2"`
		Chapters3 []struct {
			Duration string `json:"duration"`
			Image    string `json:"image"`
			URL      string `json:"url"`
		} `json:"chapters3"`
		Chapters4 []struct {
			Duration string `json:"duration"`
			Image    string `json:"image"`
			URL      string `json:"url"`
		} `json:"chapters4"`
		ValidChapterNum int    `json:"validChapterNum"`
		URL             string `json:"url"`
	} `json:"video"`
	HlsURL       string `json:"hls_url"`
	AspErrorCode string `json:"asp_error_code"`
	Manifest     struct {
		AudioMp3    string `json:"audio_mp3"`
		HlsAudioURL string `json:"hls_audio_url"`
		HlsEncURL   string `json:"hls_enc_url"`
		HlsH5EURL   string `json:"hls_h5e_url"`
		HlsEnc2URL  string `json:"hls_enc2_url"`
	} `json:"manifest"`
	ClientSid     string `json:"client_sid"`
	DefaultStream string `json:"default_stream"`
	Lc            struct {
		IspCode     string `json:"isp_code"`
		CityCode    string `json:"city_code"`
		ProviceCode string `json:"provice_code"`
		CountryCode string `json:"country_code"`
		IP          string `json:"ip"`
	} `json:"lc"`
	IsIpadSupport   string `json:"is_ipad_support"`
	Version         string `json:"version"`
	Embed           string `json:"embed"`
	IsFnMultiStream bool   `json:"is_fn_multi_stream"`
}

var vCctvComGuidReg = regexp.MustCompile(`guid = "(.*?)"`)
