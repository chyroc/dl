package parse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/chyroc/dl/pkgs/config"
	"github.com/chyroc/dl/pkgs/helper"
	"github.com/chyroc/dl/pkgs/resource"
)

func NewOpen163Com() Parser {
	return &open163Com{}
}

type open163Com struct{}

func (r *open163Com) Kind() string {
	return "open.163.com"
}

func (r *open163Com) ExampleURLs() []string {
	return []string{
		"https://open.163.com/newview/movie/free?pid=HFD3PMIPO",
		"https://open.163.com/movie/2010/6/D/6/M6TCSIN1U_M6TCSTQD6.html",
	}
}

func (r *open163Com) Parse(uri string) (resource.Resourcer, error) {
	text, err := config.ReqCli.New(http.MethodGet, uri).Text()
	if err != nil {
		return nil, err
	}
	metaCode := getMatchString(text, open163ComCharMidReg)
	metaStr, err := helper.RunJsCode(metaCode)
	if err != nil {
		return nil, err
	}
	meta := new(open163ComMeta)
	if err = json.Unmarshal([]byte(metaStr), meta); err != nil {
		return nil, err
	}

	if len(meta.State.Movie.MoiveList) == 1 {
		m := meta.State.Movie.MoiveList[0]
		return resource.NewURL(m.Title+".mp4", m.Mp4ShareURL), nil
	}

	chapters := []resource.Resourcer{}
	for _, v := range meta.State.Movie.MoiveList {
		chapters = append(chapters, resource.NewURL(fmt.Sprintf("%s_%d_%s.mp4", v.Mid, v.PNumber, v.Title), v.Mp4ShareURL))
	}

	// return download.new(meta.Data[0].Title, meta.Data[0].Title, ".mp4", chapters), err
	return resource.NewChapter(meta.Data[0].Title, chapters), err
}

var open163ComCharMidReg = regexp.MustCompile(`__NUXT__=(.*?);<`)

type open163ComMeta struct {
	Data []struct {
		Title          string `json:"title"`
		MTitle         string `json:"mTitle"`
		MDesc          string `json:"mDesc"`
		NavLink        string `json:"navLink"`
		LargeImgurl    string `json:"largeImgurl"`
		School         string `json:"school"`
		Director       string `json:"director"`
		PlayCount      int    `json:"playCount"`
		Subtitle       string `json:"subtitle"`
		Type1          string `json:"type1"`
		Desc           string `json:"desc"`
		FirstClassify  string `json:"firstClassify"`
		SecondClassify string `json:"secondClassify"`
		ThirdClassify  string `json:"thirdClassify"`
	} `json:"data"`
	State struct {
		Movie struct {
			Mid       string `json:"mid"`
			MoiveList []struct {
				Mid              string `json:"mid"`
				Plid             string `json:"plid"`
				Title            string `json:"title"`
				MLength          int    `json:"mLength"`
				ImgPath          string `json:"imgPath"`
				PNumber          int    `json:"pNumber"`
				Subtitle         string `json:"subtitle"`
				SubtitleLanguage string `json:"subtitleLanguage"`
				WebURL           string `json:"webUrl"`
				ShortWebURL      string `json:"shortWebUrl"`
				Hits             int    `json:"hits"`
				Mp4HdSize        int    `json:"mp4HdSize"`
				Mp4HdURL         string `json:"mp4HdUrl"`
				M3U8HdURL        string `json:"m3u8HdUrl"`
				Mp4ShareURL      string `json:"mp4ShareUrl"`
				Description      string `json:"description"`
			} `json:"moiveList"`
		} `json:"movie"`
	} `json:"state"`
}
