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

func NewHaokanBaiduCom() Parser {
	return &haokanBaiduCom{}
}

type haokanBaiduCom struct{}

func (r *haokanBaiduCom) Kind() string {
	return "haokan.baidu.com"
}

func (r *haokanBaiduCom) Parse(uri string) (resource.Resource, error) {
	title, videoURL, err := r.getMeta(uri)
	if err != nil {
		return nil, err
	}

	return resource.NewURL(title+".mp4", videoURL), nil
}

func (r *haokanBaiduCom) getMeta(uri string) (string, string, error) {
	text, err := config.ReqCli.New(http.MethodGet, uri).Text()
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
