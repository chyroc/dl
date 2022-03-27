package parse

import (
	"net/http"
	"regexp"

	"github.com/chyroc/dl/internal/config"
	"github.com/chyroc/dl/internal/resource"
)

func NewMSohuCom() Parser {
	return &mSohuCom{}
}

type mSohuCom struct{}

func (r *mSohuCom) Kind() string {
	return "m.sohu.com"
}

func (r *mSohuCom) ExampleURLs() []string {
	return []string{"https://m.sohu.com/a/490513509_120538293"}
}

func (r *mSohuCom) Parse(uri string) (resource.Resource, error) {
	text, err := config.ReqCli.New(http.MethodGet, uri).Text()
	if err != nil {
		return nil, err
	}
	title := getMatchString(text, mSohuComTitleReg) + "_" + getMatchString(text, mSohuComAuthorReg)
	url := getMatchString(text, mSohuComUrlReg)

	return resource.NewURL(title+".mp4", url), nil
}

var (
	mSohuComUrlReg    = regexp.MustCompile(`data-url="(.*?)"`)
	mSohuComTitleReg  = regexp.MustCompile(`title: '(.*?)'`)
	mSohuComAuthorReg = regexp.MustCompile(`authorId: '(.*?)'`)
)
