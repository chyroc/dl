package identify

import (
	"fmt"
	"net/url"

	"github.com/chyroc/dl/internal/parse"
)

func Identify(uri string) (parse.Parser, error) {
	uriParsed, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("parse url: %q failed: %w", uri, err)
	}
	switch uriParsed.Host {
	case "video.sina.com.cn":
		return parse.NewVideoSinaComCn(), nil
	case "haokan.baidu.com":
		return parse.NewHaokanBaiduCom(), nil
	case "v.youku.com":
		return parse.NewYoukuCom(), nil
	}

	return nil, nil
}
