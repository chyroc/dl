package identify

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/chyroc/dl/internal/parse"
)

func Identify(uri string) (parse.Parser, error) {
	uriParsed, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("parse url: %q failed: %w", uri, err)
	}
	parser, ok := hostToParserRegister[uriParsed.Host]
	if ok {
		return parser, nil
	}

	parser, ok = fileTypeToParserRegister[strings.ToLower(filepath.Ext(uri))]
	if ok {
		return parser, nil
	}

	return nil, nil
}

var (
	hostToParserRegister     = map[string]parse.Parser{}
	fileTypeToParserRegister = map[string]parse.Parser{}
	ExampleURLs              []string
)

func register(parser parse.Parser) {
	kind := parser.Kind()
	if strings.HasPrefix(kind, "filetype.") {
		fileTypeToParserRegister[strings.ToLower(kind[len("filetype"):])] = parser
	} else {
		for _, v := range strings.Split(kind, ",") {
			hostToParserRegister[v] = parser
		}
	}

	ExampleURLs = append(ExampleURLs, parser.ExampleURLs()...)
}
