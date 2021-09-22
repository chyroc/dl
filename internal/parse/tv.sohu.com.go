package parse

import (
	"fmt"
	"net/http"
	"regexp"
	"sync"

	"github.com/chyroc/dl/internal/config"
	"github.com/chyroc/dl/internal/download"
)

func NewTvSohuCom() Parser {
	return &tvSohuCom{}
}

type tvSohuCom struct{}

func (r *tvSohuCom) Kind() string {
	return "tv.sohu.com"
}

func (r *tvSohuCom) Parse(uri string) (download.Downloader, error) {
	htmlMeta, err := r.getHTMLMeta(uri)
	if err != nil {
		return nil, err
	}

	videoMeta, err := r.getVideoClips(uri, htmlMeta.Vid)
	if err != nil {
		return nil, err
	}

	urls := make([]string, len(videoMeta.Data.Su))
	wg := new(sync.WaitGroup)
	var finalErr error
	for i, v := range videoMeta.Data.Su {
		wg.Add(1)
		go func(i int, v string) {
			defer wg.Done()

			if finalErr != nil {
				return
			}

			url, err := r.getVideoURL(v)
			if err != nil {
				finalErr = err
				return
			}
			urls[i] = url
		}(i, v)
	}
	wg.Wait()

	title := fmt.Sprintf("%s_%d", videoMeta.Data.TvName, videoMeta.Tvid)
	return download.NewDownloadURLList(title, title+".mp4", urls), nil
}

func (r *tvSohuCom) getVideoClips(originURL, vid string) (*tvSohuComGetVideoClipsResp, error) {
	uri := fmt.Sprintf("https://hot.vrs.sohu.com/vrs_flash.action?vid=%s&ver=1&ssl=1&pflag=pch5", vid)
	resp := new(tvSohuComGetVideoClipsResp)
	err := config.ReqCli.New(http.MethodGet, uri).WithHeaders(prepareCommonHeader(originURL, map[string]string{
		"accept-language": "zh-CN,zh;q=0.9,en;q=0.8",
	})).Unmarshal(resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *tvSohuCom) getHTMLMeta(uri string) (*tvSohuComHtmlMeta, error) {
	text, err := config.ReqCli.New(http.MethodGet, uri).Text()
	if err != nil {
		return nil, err
	}
	vid := getMatchString(text, tvSohuComVidReg)

	return &tvSohuComHtmlMeta{
		Vid: vid,
	}, nil
}

func (r *tvSohuCom) getVideoURL(key string) (string, error) {
	uri := fmt.Sprintf("https://data.vod.itc.cn/ip?new=%s&num=1&ch=tv&pt=1&pg=2&prod=h5n", key)
	resp := new(tvSohuComGetVideoURLResp)

	err := config.ReqCli.New(http.MethodGet, uri).Unmarshal(resp)
	if err != nil {
		return "", err
	}
	return resp.Servers[0].URL, nil
}

type tvSohuComHtmlMeta struct {
	Vid string
}

type tvSohuComGetVideoClipsResp struct {
	URL   string `json:"url"`
	Pvpic struct {
		Big   string `json:"big"`
		Small string `json:"small"`
	} `json:"pvpic"`
	Tvid int   `json:"tvid"`
	Syst int64 `json:"syst"`
	Data struct {
		TvName        string    `json:"tvName"`
		SubName       string    `json:"subName"`
		Ch            string    `json:"ch"`
		Fps           int       `json:"fps"`
		IPLimit       int       `json:"ipLimit"`
		Width         int       `json:"width"`
		ClipsURL      []string  `json:"clipsURL"`
		Version       int       `json:"version"`
		ClipsBytes    []int     `json:"clipsBytes"`
		Num           int       `json:"num"`
		CoverImg      string    `json:"coverImg"`
		Height        int       `json:"height"`
		TotalDuration float64   `json:"totalDuration"`
		TotalBytes    int       `json:"totalBytes"`
		ClipsDuration []float64 `json:"clipsDuration"`
		Orifee        int       `json:"orifee"`
		Ck            []string  `json:"ck"`
		Hc            []string  `json:"hc"`
		Su            []string  `json:"su"`
	} `json:"data"`
	Keyword string `json:"keyword"`
	Cmscat  string `json:"cmscat"`
}

type tvSohuComGetVideoURLResp struct {
	Servers []struct {
		Nid   int    `json:"nid"`
		Isp2P int    `json:"isp2p"`
		URL   string `json:"url"`
	} `json:"servers"`
}

var tvSohuComVidReg = regexp.MustCompile(`vid="(\d+)"`)
