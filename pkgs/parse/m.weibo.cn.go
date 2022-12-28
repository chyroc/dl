package parse

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/chyroc/dl/pkgs/config"
	"github.com/chyroc/dl/pkgs/resource"
)

func NewMobileWeiboCn() Parser {
	return &mobileWeiboCn{}
}

type mobileWeiboCn struct{}

func (r *mobileWeiboCn) Kind() string {
	return "m.weibo.cn"
}

func (r *mobileWeiboCn) ExampleURLs() []string {
	return []string{""}
}

func (r *mobileWeiboCn) Parse(uri string) (resource.Resourcer, error) {
	header := prepareCommonHeader(uri, nil)
	text, err := config.ReqCli.New(http.MethodGet, uri).WithHeaders(header).Text()
	if err != nil {
		return nil, err
	}

	title := getMatchString(text, mWeiboCnTitleReg)
	url := getMatchString(text, mWeiboCnUrlReg)
	user := getMatchString(text, mWeiboCnUserReg)
	id := getMatchString(uri, mWeiboCnIdReg)

	return resource.NewURL(fmt.Sprintf("%s_%s_%s.mp4", title, user, id), url), nil
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
