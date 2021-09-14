package parse

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/chyroc/dl/internal/config"
	"github.com/chyroc/dl/internal/download"
	"github.com/chyroc/dl/internal/helper"
	"github.com/chyroc/gorequests"
)

func NewVideoSinaComCn() Parser {
	return &videoSinaComCn{}
}

type videoSinaComCn struct{}

func (r *videoSinaComCn) Parse(uri string) (download.Downloader, error) {
	videoID, err := r.getVideoID(uri)
	if err != nil {
		return nil, err
	}
	return r.getVideoMeta(uri, videoID)
}

func (r *videoSinaComCn) getVideoID(uri string) (int64, error) {
	text, err := gorequests.New(http.MethodGet, uri).WithLogger(config.WithLogger()).Text()
	if err != nil {
		return 0, err
	}
	match := videoSinaComCnVideoIDReg.FindStringSubmatch(text)
	if len(match) == 2 {
		videoID, _ := strconv.ParseInt(match[1], 10, 64)
		return videoID, nil
	}

	return 0, fmt.Errorf("parse %q video_id failed", uri)
}

func (r *videoSinaComCn) getVideoMeta(originURL string, videoID int64) (download.Downloader, error) {
	uri := "http://api.ivideo.sina.com.cn/public/video/play"
	query := map[string]string{
		"video_id": strconv.FormatInt(videoID, 10),
		"appver":   "V11220.210521.02",
		"appname":  "sinaplayer_pc",
		"applt":    "web",
		"tags":     "sinaplayer_pc",
		"player":   "all",
	}
	header := map[string]string{
		"Host":       "api.ivideo.sina.com.cn",
		"Connection": "keep-alive",
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36",
		"Accept":     "*/*",
		"Origin":     "http://video.sina.com.cn",
		"Referer":    originURL,
		// "Accept-Encoding": "gzip, deflate",
		"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
	}
	resp := new(videoSinaComCnGetVideoMetaResp)

	err := gorequests.New(http.MethodGet, uri).WithQuerys(query).WithHeaders(header).WithLogger(config.WithLogger()).Unmarshal(resp)
	if err != nil {
		return nil, err
	} else if resp.Code != 1 {
		return nil, fmt.Errorf(resp.Message)
	}

	pkgs := []*download.Pkg{}
	for _, v := range resp.Data.Videos {
		pkgs = append(pkgs, &download.Pkg{
			Size:       helper.MayStringToInt64(v.Size),
			Definition: download.MayConvertDefinition(v.Definition),
			URL:        v.DispatchResult.URL,
			Type:       v.Type,
		})
	}

	return download.NewDownloadURL(resp.Data.Title, pkgs), nil
}

type videoSinaComCnGetVideoMetaResp struct {
	Message string `json:"Message"`
	Code    int    `json:"code"`
	Data    struct {
		CreateTime string `json:"create_time"`
		Image      string `json:"image"`
		Length     string `json:"length"`
		Title      string `json:"title"`
		Videos     []struct {
			Codec          string `json:"codec"`
			Definition     string `json:"definition"`
			FileID         string `json:"file_id"`
			Height         string `json:"height"`
			Length         string `json:"length"`
			Md5            string `json:"md5"`
			Size           string `json:"size"`
			Status         string `json:"status"`
			Type           string `json:"type"`
			Width          string `json:"width"`
			Avc            string `json:"avc"`
			DispatchResult struct {
				Result string `json:"result"`
				URL    string `json:"url"`
				Bakurl string `json:"bakurl"`
			} `json:"dispatch_result"`
		} `json:"videos"`
	} `json:"data"`
	Error        string `json:"error"`
	ErrorMessage string `json:"errorMessage"`
}

var videoSinaComCnVideoIDReg = regexp.MustCompile(`video_id:'?(\d+)'?,`)
