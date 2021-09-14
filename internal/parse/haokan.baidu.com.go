package parse

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chyroc/dl/internal/config"
	"github.com/chyroc/dl/internal/download"
	"github.com/chyroc/gorequests"
)

func NewHaokanBaiduCom() Parser {
	return &haokanBaiduCom{}
}

type haokanBaiduCom struct{}

func (r *haokanBaiduCom) Kind() string {
	return "haokan.baidu.com"
}

func (r *haokanBaiduCom) Parse(uri string) (download.Downloader, error) {
	title, videoURL, err := r.getMeta(uri)
	if err != nil {
		return nil, err
	}

	spec := &download.Specification{
		Size:       0,
		Definition: "",
		URL:        videoURL,
	}

	return download.NewDownloadURL(title, title+".mp4", []*download.Specification{spec}), nil
}

func (r *haokanBaiduCom) getMeta(uri string) (string, string, error) {
	text, err := gorequests.New(http.MethodGet, uri).WithLogger(config.WithLogger()).Text()
	if err != nil {
		return "", "", err
	}
	videoURL := ""
	match := haokanBaiduComURLReg.FindStringSubmatch(text)
	if len(match) == 2 {
		videoURL = match[1]
		if strings.Contains(videoURL, "\\/") {
			videoURL = strings.ReplaceAll(videoURL, "\\/", "/")
		}
	} else {
		return "", "", fmt.Errorf("parse %q get video url failed", uri)
	}
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(text))
	title := doc.Find("title").Text()

	return title, videoURL, nil
}

var haokanBaiduComURLReg = regexp.MustCompile(`"playurl":"(.*?)"`)
