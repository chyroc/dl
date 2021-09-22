package identify

import (
	"fmt"
	"net/url"
	"strings"

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
	kind := parser.Kind()
	for _, v := range strings.Split(kind, ",") {
		parserRegister[v] = parser
	}
}

func init() {
	register(parse.NewVideoSinaComCn())
	register(parse.NewHaokanBaiduCom())
	register(parse.NewVYoukuCom())
	register(parse.NewMobileWeiboCn())
	register(parse.NewVideoCaixinCom())
	register(parse.NewTvSohuCom())
	register(parse.NewVCctvCom())
	register(parse.NewA36krCom())
	register(parse.NewMSohuCom())
	register(parse.NewWwwZhihuCom())
	register(parse.NewVDouyinCom())
	register(parse.NewOpen163Com())
	register(parse.NewWwwMissevanCom())
	register(parse.NewWww333tttCom())
}
