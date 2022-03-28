package parse

import (
	"net/http"
	"regexp"

	"github.com/chyroc/dl/pkgs/config"
	"github.com/chyroc/dl/pkgs/resource"
)

func NewWww333tttCom() Parser {
	return &www333tttCom{}
}

type www333tttCom struct{}

func (r *www333tttCom) Kind() string {
	return "www.333ttt.com"
}

func (r *www333tttCom) ExampleURLs() []string {
	return []string{"http://www.333ttt.com/up/yy6182865.html"}
}

func (r *www333tttCom) Parse(uri string) (resource.Resource, error) {
	text, err := config.ReqCli.New(http.MethodGet, uri).Text()
	if err != nil {
		return nil, err
	}

	title := getMatchString(text, www333tttComNameReg)
	url := getMatchString(text, www333tttComUrlReg) + ".mp3"
	return resource.NewURL(title+".mp3", url), nil
}

var (
	www333tttComNameReg = regexp.MustCompile(`id="play_name">(.*?)</span>`)
	www333tttComUrlReg  = regexp.MustCompile(`href="(.*?)\.mp3"`)
)
