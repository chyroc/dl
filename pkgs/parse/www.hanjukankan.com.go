package parse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/chyroc/dl/pkgs/config"
	"github.com/chyroc/dl/pkgs/helper"
	"github.com/chyroc/dl/pkgs/resource"
)

func NewWwwHanjukankanCom() Parser {
	return &wwwHanjukankanCom{}
}

type wwwHanjukankanCom struct{}

func (r *wwwHanjukankanCom) Kind() string {
	return "www.hanjukankan.com"
}

func (r *wwwHanjukankanCom) ExampleURLs() []string {
	return []string{
		"https://www.hanjukankan.com/movie/index1254.html",
		"https://www.hanjukankan.com/play/1254-0-1.html",
	}
}

func (r *wwwHanjukankanCom) Parse(uri string) (resource.Resourcer, error) {
	vid := ""
	fetchAll := true
	if match := vidReg.FindStringSubmatch(uri); len(match) != 2 {
		match = vidOneReg.FindStringSubmatch(uri)
		if len(match) != 2 {
			return nil, fmt.Errorf("parse uri fail: %s", uri)
		}
		fetchAll = false
		vid = match[1]
	} else {
		fetchAll = true
		vid = match[1]
	}

	if fetchAll {
		return r.fetchAll(uri, vid)
	} else {
		return r.fetchOne(uri, vid)
	}
}

func (r *wwwHanjukankanCom) fetchAll(uri, vid string) (resource.Resourcer, error) {
	title, err := r.getTitle(uri)
	if err != nil {
		return nil, err
	}
	urls, err := r.getURLs(vid)
	if err != nil {
		return nil, err
	}
	tmp := []resource.Resourcer{}
	for idx, v := range urls {
		chapterResource, err := resource.NewM3U8(100, fmt.Sprintf("%s-%d.mp4", title, idx+1), v)
		if err != nil {
			return nil, err
		}
		tmp = append(tmp, chapterResource)
	}
	return resource.NewChapter(title, tmp), nil
}

func (r *wwwHanjukankanCom) fetchOne(uri, vid string) (resource.Resourcer, error) {
	match := chapterIndexReg.FindStringSubmatch(uri)
	if len(match) != 2 {
		return nil, fmt.Errorf("parse uri fail: %s", uri)
	}
	index, _ := strconv.Atoi(match[1])
	title, err := r.getTitle(uri)
	if err != nil {
		return nil, err
	}
	urls, err := r.getURLs(vid)
	if err != nil {
		return nil, err
	}
	if index >= len(urls) {
		return nil, fmt.Errorf("index out of range: %d", index)
	}
	url := urls[index]
	return resource.NewM3U8(100, fmt.Sprintf("%s.mp4", title), url)
}

func (r *wwwHanjukankanCom) getTitle(uri string) (string, error) {
	text, err := config.ReqCli.New(http.MethodGet, uri).Text()
	if err != nil {
		return "", fmt.Errorf("get html of %s fail: %s", uri, err)
	}

	return helper.HtmlTitle(text), nil
}

func (r *wwwHanjukankanCom) getURLs(vid string) ([]string, error) {
	uri := fmt.Sprintf("https://www.hanjukankan.com/ass.php?vid=%s&vfrom=0", vid)
	text, err := config.ReqCli.New(http.MethodGet, uri).Text()
	if err != nil {
		return nil, fmt.Errorf("get json of %s fail: %s", uri, err)
	}
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "(") {
		text = text[1:]
	}
	if strings.HasSuffix(text, ");") {
		text = text[:len(text)-2]
	}

	res := new(parseVideoURLResp)
	err = json.Unmarshal([]byte(text), res)
	if err != nil {
		return nil, fmt.Errorf("parse json of %q fail: %s", text, err)
	}

	return res.S.Video, nil
}

var vidOneReg = regexp.MustCompile(`/play/(\d+)-.*?.html`)

// /movie/index1254.html
var vidReg = regexp.MustCompile(`/movie/index(\d+).html`)

// https://www.hanjukankan.com/play/1254-0-10.html
var chapterIndexReg = regexp.MustCompile(`/play/\d+-\d+-(\d+).html`)

type parseVideoURLResp struct {
	S struct {
		Video []string `json:"video"`
	} `json:"s"`
}
