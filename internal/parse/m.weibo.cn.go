package parse

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/chyroc/dl/internal/config"
	"github.com/chyroc/dl/internal/download"
)

func NewMobileWeiboCn() Parser {
	return &mobileWeiboCn{}
}

type mobileWeiboCn struct{}

func (r *mobileWeiboCn) Kind() string {
	return "m.weibo.cn"
}

func (r *mobileWeiboCn) Parse(uri string) (download.Downloader, error) {
	header := prepareCommonHeader(uri, nil)
	text, err := config.ReqCli.New(http.MethodGet, uri).WithHeaders(header).Text()
	if err != nil {
		return nil, err
	}

	title := getMatchString(text, mWeiboCnTitleReg)
	url := getMatchString(text, mWeiboCnUrlReg)
	user := getMatchString(text, mWeiboCnUserReg)
	id := getMatchString(uri, mWeiboCnIdReg)
	fmt.Println(uri, mWeiboCnIdReg, mWeiboCnIdReg.FindStringSubmatch(uri))
	specs := []*download.Specification{
		{
			Size:       0,
			Definition: download.DefinitionHD,
			URL:        url,
		},
	}

	return download.NewDownloadURL(title, fmt.Sprintf("%s_%s_%s.mp4", title, user, id), false, specs), nil
}

var (
	mWeiboCnTitleReg = regexp.MustCompile(`"title": "(.*?)"`)
	mWeiboCnUrlReg   = regexp.MustCompile(`"stream_url_hd": "(.*?)"`)
	mWeiboCnUserReg  = regexp.MustCompile(`"screen_name": "(.*?)"`)
	mWeiboCnIdReg    = regexp.MustCompile(`m\.weibo\.cn/detail/(\d+)`)
)

func getMatchString(s string, reg *regexp.Regexp) string {
	match := reg.FindStringSubmatch(s)
	if len(match) == 2 {
		return match[1]
	}
	return ""
}

func getMatchStringByRegs(s string, regs []*regexp.Regexp) string {
	for _, reg := range regs {
		if v := getMatchString(s, reg); v != "" {
			return v
		}
	}
	return ""
}
