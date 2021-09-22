package parse

import (
	"net/http"
	"regexp"

	"github.com/chyroc/dl/internal/config"
	"github.com/chyroc/dl/internal/download"
)

func NewWww333tttCom() Parser {
	return &www333tttCom{}
}

type www333tttCom struct{}

func (r *www333tttCom) Kind() string {
	return "www.333ttt.com"
}

func (r *www333tttCom) Parse(uri string) (download.Downloader, error) {
	text, err := config.ReqCli.New(http.MethodGet, uri).Text()
	if err != nil {
		return nil, err
	}

	title := getMatchString(text, www333tttComNameReg)
	url := getMatchString(text, www333tttComUrlReg) + ".mp3"
	return download.NewDownloadURL(title, title+".mp3", false, []*download.Specification{{URL: url}}), nil
}

var (
	www333tttComNameReg = regexp.MustCompile(`id="play_name">(.*?)</span>`)
	www333tttComUrlReg  = regexp.MustCompile(`href="(.*?)\.mp3"`)
)
