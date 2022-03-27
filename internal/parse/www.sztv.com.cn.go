package parse

import (
	"github.com/chyroc/dl/internal/resource"
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

func (r *wwwSztvComCn) Parse(uri string) (resource.Resource, error) {
	panic("")
	// text, err := config.ReqCli.New(http.MethodGet, uri).Text()
	// if err != nil {
	// 	return nil, err
	// }
	// d, _ := goquery.NewDocumentFromReader(strings.NewReader(text))
	// d.Find(".videoInfo").Each(func(i int, s *goquery.Selection) {
	// 	title := s.AttrOr("title", "")
	// 	src := s.AttrOr("src", "")
	// 	if src == "" {
	// 		return
	// 	}
	//
	// 	download.NewDownloadBggee()
	// })
	//
	// return download.NewDownloadURL(title, title+".mp4", false, []*download.Specification{spec}), nil
}
