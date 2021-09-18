package parse

import (
	"net/http"
	"regexp"

	"github.com/chyroc/dl/internal/config"
	"github.com/chyroc/dl/internal/download"
)

func NewMSohuCom() Parser {
	return &mSohuCom{}
}

type mSohuCom struct{}

func (r *mSohuCom) Kind() string {
	return "m.sohu.com"
}

func (r *mSohuCom) Parse(uri string) (download.Downloader, error) {
	text, err := config.ReqCli.New(http.MethodGet, uri).Text()
	if err != nil {
		return nil, err
	}
	title := getMatchString(text, mSohuComTitleReg) + "_" + getMatchString(text, mSohuComAuthorReg)
	url := getMatchString(text, mSohuComUrlReg)
	return download.NewDownloadURL(title, title+".mp4", false, []*download.Specification{{
		URL: url,
	}}), nil
}

var (
	mSohuComUrlReg    = regexp.MustCompile(`data-url="(.*?)"`)
	mSohuComTitleReg  = regexp.MustCompile(`title: '(.*?)'`)
	mSohuComAuthorReg = regexp.MustCompile(`authorId: '(.*?)'`)
)
