package parse

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/chyroc/dl/pkgs/config"
	"github.com/chyroc/dl/pkgs/resource"
)

func NewNewsCctvCom() Parser {
	return &newsCctvCom{}
}

type newsCctvCom struct{}

func (r *newsCctvCom) Kind() string {
	return "news.cctv.com"
}

func (r *newsCctvCom) ExampleURLs() []string {
	return []string{"http://news.cctv.com/2021/10/02/VIDEjj6VR17N4AEAIrUXmwWV211002.shtml"}
}

func (r *newsCctvCom) Parse(uri string) (resource.Resourcer, error) {
	guid, err := r.getGuid(uri)
	if err != nil {
		return nil, err
	}
	info, err := r.getVideoInfo(guid)
	if err != nil {
		return nil, err
	}
	downloadURL, err := r.getDownloadURL(info.HlsURL)

	return resource.NewURL(info.Title+".mp4", downloadURL), nil
}

func (r *newsCctvCom) getGuid(uri string) (string, error) {
	text, err := config.ReqCli.New(http.MethodGet, uri).Text()
	if err != nil {
		return "", err
	}
	return getMatchString(text, newsCctvComGuidRegex), nil
}

func (r *newsCctvCom) getVideoInfo(guid string) (*videoInfo, error) {
	resp := new(videoInfo)
	err := config.ReqCli.New(http.MethodGet, "http://vdn.apps.cntv.cn/api/getHttpVideoInfo.do?pid="+guid).Unmarshal(resp)
	return resp, err
}

func (r *newsCctvCom) getDownloadURL(hlsURL string) (string, error) {
	host, _ := url.Parse(hlsURL)
	text, err := config.ReqCli.New(http.MethodGet, hlsURL).Text()
	if err != nil {
		return "", err
	}
	list := []string{}
	for _, v := range strings.Split(text, "\n") {
		if !strings.HasPrefix(v, "#") && strings.HasSuffix(v, ".m3u8") {
			list = append(list, v)
		}
	}
	sort.Slice(list, func(i, j int) bool {
		ii := strings.Split(filepath.Base(list[i]), "/")
		jj := strings.Split(filepath.Base(list[j]), "/")
		iii, _ := strconv.Atoi(ii[len(ii)-1])
		jjj, _ := strconv.Atoi(jj[len(jj)-1])
		return iii > jjj
	})
	return fmt.Sprintf("%s://%s%s", host.Scheme, host.Host, list[0]), nil
}

var newsCctvComGuidRegex = regexp.MustCompile(`guid = "(.*?)";`)

type videoInfo struct {
	Title  string `json:"title"`
	HlsURL string `json:"hls_url"`
}
