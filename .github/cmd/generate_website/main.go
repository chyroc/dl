package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"text/template"

	"github.com/chyroc/dl/.github/cmd/cmd_helper"
)

func main() {
	URL, err := url.Parse(os.Args[1])
	cmd_helper.Assert(err)

	uri := URL.Host
	req := &generateGoCodeReq{
		Host:               uri,
		LowerCamelCaseHost: hostToLowerCamelCase(uri),
		TitleCamelCaseHost: hostToTitleCamelCase(uri),
	}

	code := generateGoCode(req)

	cmd_helper.Assert(ioutil.WriteFile(fmt.Sprintf("./pkgs/parse/%s.go", uri), []byte(code), 0o666))

	code = generateGoTestCode(req)
	cmd_helper.Assert(ioutil.WriteFile(fmt.Sprintf("./pkgs/parse/%s_test.go", uri), []byte(code), 0o666))
}

func hostToTitleCamelCase(s string) string {
	s = hostToLowerCamelCase(s)
	return string(append([]rune{[]rune(s)[0] + 'A' - 'a'}, []rune(s)[1:]...))
}

func hostToLowerCamelCase(s string) string {
	res := []rune{}
	bigger := false
	for i, v := range []rune(s) {
		if i == 0 {
			res = append(res, v)
		} else if v == '.' {
			bigger = true
		} else if bigger {
			if v >= 'a' && v <= 'z' {
				res = append(res, v+'A'-'a')
			} else {
				res = append(res, v)
			}
			bigger = false
		} else {
			res = append(res, v)
		}
	}
	if res[0] >= '0' && res[0] <= '9' {
		res = append([]rune{'a'}, res...)
	}
	return string(res)
}

type generateGoCodeReq struct {
	Host               string
	LowerCamelCaseHost string
	TitleCamelCaseHost string
}

func generateGoCode(req *generateGoCodeReq) string {
	t, err := template.New("").Parse(goTemplate)
	cmd_helper.Assert(err)
	buf := new(bytes.Buffer)
	err = t.Execute(buf, req)
	cmd_helper.Assert(err)
	return buf.String()
}

var goTemplate = `package parse

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/chyroc/dl/pkgs/resource"
)

func New{{ .TitleCamelCaseHost }}() Parser {
	return &{{ .LowerCamelCaseHost }}{}
}

type {{ .LowerCamelCaseHost }} struct{}

func (r *{{ .LowerCamelCaseHost }}) Kind() string {
	return "{{ .Host }}"
}

func (r *{{ .LowerCamelCaseHost }}) ExampleURLs() []string {
	return []string{}
}

func (r *{{ .LowerCamelCaseHost }}) Parse(uri string) (resource.Resourcer, error) {
	panic("")
}
`

func generateGoTestCode(req *generateGoCodeReq) string {
	t, err := template.New("").Parse(goTestTemplate)
	cmd_helper.Assert(err)
	buf := new(bytes.Buffer)
	err = t.Execute(buf, req)
	cmd_helper.Assert(err)
	return buf.String()
}

var goTestTemplate = `package parse_test

import (
	"testing"

	"github.com/chyroc/dl/pkgs/parse"
	"github.com/chyroc/go-assert"
)

func Test_{{ .LowerCamelCaseHost }}(t *testing.T) {
	as := assert.New(t, assert.WithFailRerun(5))

	{
		res, err := parse.New{{ .TitleCamelCaseHost }}().Parse("...")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("...", res.Title())
	}
}

`
