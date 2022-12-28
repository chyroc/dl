package parse

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chyroc/dl/pkgs/config"
	"github.com/chyroc/dl/pkgs/resource"
)

func NewWwwSztvComCn() Parser {
	return &wwwSztvComCn{}
}

type wwwSztvComCn struct{}

func (r *wwwSztvComCn) Kind() string {
	return "www.sztv.com.cn"
}

func (r *wwwSztvComCn) ExampleURLs() []string {
	return []string{""}
}

func (r *wwwSztvComCn) Parse(uri string) (resource.Resourcer, error) {
	text, err := config.ReqCli.New(http.MethodGet, uri).Text()
	if err != nil {
		return nil, err
	}
	d, _ := goquery.NewDocumentFromReader(strings.NewReader(text))

	pageTitle := strings.TrimSpace(d.Find("title").Text())
	chapterList := []resource.Resourcer{}
	d.Find(".videoInfo").Each(func(i int, s *goquery.Selection) {
		title := s.AttrOr("title", "")
		src := s.AttrOr("src", "")
		if src == "" {
			return
		}

		chapterList = append(chapterList, resource.NewURL(title+".mp4", src))
	})

	return resource.NewURLChapter(pageTitle, chapterList), nil
}
