package parse

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/chyroc/dl/internal/config"
	"github.com/chyroc/dl/internal/resource"
)

func NewA36krCom() Parser {
	return &a36krCom{}
}

type a36krCom struct{}

func (r *a36krCom) Kind() string {
	return "36kr.com"
}

func (r *a36krCom) ExampleURLs() []string {
	return []string{"https://36kr.com/video/1402287552562048"}
}

func (r *a36krCom) Parse(uri string) (resource.Resource, error) {
	text, err := config.ReqCli.New(http.MethodGet, uri).Text()
	if err != nil {
		return nil, err
	}

	info := getMatchString(text, a36krComInfoReg)

	resp := new(a36krComInfoGetResp)
	if err = json.Unmarshal([]byte(info), resp); err != nil {
		return nil, err
	}
	data := resp.VideoDetail.Data

	return resource.NewURL(data.WidgetTitle+".mp4", data.URL), nil
}

var a36krComInfoReg = regexp.MustCompile(`<script>window.initialState=(.*?)</script>`)

type a36krComInfoGetResp struct {
	VideoDetail struct {
		Code int              `json:"code"`
		Data a36krComInfoMeta `json:"data"`
	} `json:"videoDetail"`
}

type a36krComInfoMeta struct {
	ItemID         int64  `json:"itemId"`
	WidgetTitle    string `json:"widgetTitle"`
	WidgetContent  string `json:"widgetContent"`
	WidgetImage    string `json:"widgetImage"`
	Duration       int    `json:"duration"`
	PublishTime    int64  `json:"publishTime"`
	AuthorName     string `json:"authorName"`
	AuthorRoute    string `json:"authorRoute"`
	AuthorFace     string `json:"authorFace"`
	URL            string `json:"url"`
	Filesize       int64  `json:"filesize"`
	URL384         string `json:"url384"`
	Filesize384    int    `json:"filesize384"`
	URL512         string `json:"url512"`
	Filesize512    int    `json:"filesize512"`
	URL1152        string `json:"url1152"`
	Filesize1152   int    `json:"filesize1152"`
	DefinitionBean struct {
		VideoID        int64  `json:"videoId"`
		URLOrigin      string `json:"urlOrigin"`
		FileSizeOrigin int    `json:"fileSizeOrigin"`
		URLFHD         string `json:"urlFHD"`
		FileSizeFHD    int    `json:"fileSizeFHD"`
		URLHD          string `json:"urlHD"`
		FileSizeHD     int    `json:"fileSizeHD"`
		URLSD          string `json:"urlSD"`
		FileSizeSD     int    `json:"fileSizeSD"`
		URLLD          string `json:"urlLD"`
		FileSizeLD     int    `json:"fileSizeLD"`
	} `json:"definitionBean"`
	AuthorSummary string `json:"authorSummary"`
	ShapeType     int    `json:"shapeType"`
}
