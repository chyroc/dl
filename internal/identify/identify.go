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
	if uriParsed.Host == "video.sina.com.cn" {
		return parse.NewVideoSinaComCn(), nil
	}

	return nil, nil
}
