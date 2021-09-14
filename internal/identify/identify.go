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
	parser, ok := parserRegister[uriParsed.Host]
	if ok {
		return parser, nil
	}

	return nil, nil
}

var parserRegister = map[string]parse.Parser{}

func register(parser parse.Parser) {
	parserRegister[parser.Kind()] = parser
}

func init() {
	register(parse.NewVideoSinaComCn())
	register(parse.NewHaokanBaiduCom())
	register(parse.NewVYoukuCom())
	register(parse.NewMobileWeiboCn())
}
