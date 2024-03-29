package parse

import (
	"github.com/chyroc/dl/pkgs/resource"
)

func NewYQqCom() Parser {
	return &yQqCom{}
}

type yQqCom struct{}

func (r *yQqCom) Kind() string {
	return "y.qq.com"
}

func (r *yQqCom) ExampleURLs() []string {
	return []string{"https://y.qq.com/n/yqq/song/002Zkt5S2z8JZx.html"}
}

func (r *yQqCom) Parse(uri string) (resource.Resourcer, error) {
	panic("")
	// req, err := tencent.Parse(uri)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// if err = req.Do(); err != nil {
	// 	return nil, err
	// }
	//
	// mp3List, err := req.Extract()
	// if err != nil {
	// 	return nil, err
	// }
	//
	// if len(mp3List) == 0 {
	// 	return nil, fmt.Errorf("find no mp3")
	// }
	// if len(mp3List) == 1 {
	// 	return download.NewDownloadMp3(mp3List[0], "", false), nil
	// }
	// return download.NewDownloadMp3List(mp3List[0].SavePath, mp3List[0].SavePath, mp3List), nil
}
