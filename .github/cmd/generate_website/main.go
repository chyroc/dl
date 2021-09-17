package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
)

func main() {
	uri := os.Args[1]
	req := &generateGoCodeReq{
		Host:               uri,
		LowerCamelCaseHost: hostToLowerCamelCase(uri),
		TitleCamelCaseHost: hostToTitleCamelCase(uri),
	}

	code := generateGoCode(req)

	assert(ioutil.WriteFile(fmt.Sprintf("./internal/parse/%s.go", uri), []byte(code), 0o666))

	code = generateGoTestCode(req)
	assert(ioutil.WriteFile(fmt.Sprintf("./internal/parse/%s_test.go", uri), []byte(code), 0o666))
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
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
			res = append(res, v+'A'-'a')
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
	assert(err)
	buf := new(bytes.Buffer)
	err = t.Execute(buf, req)
	assert(err)
	return buf.String()
}

var goTemplate = `package parse

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/chyroc/dl/internal/config"
	"github.com/chyroc/dl/internal/download"
)

func New{{ .TitleCamelCaseHost }}() Parser {
	return &{{ .LowerCamelCaseHost }}{}
}

type {{ .LowerCamelCaseHost }} struct{}

func (r *{{ .LowerCamelCaseHost }}) Kind() string {
	return "{{ .Host }}"
}

func (r *{{ .LowerCamelCaseHost }}) Parse(uri string) (download.Downloader, error) {

	return download.NewDownloadURL(title, title+".mp4", []*download.Specification{spec}), nil
}

`

func generateGoTestCode(req *generateGoCodeReq) string {
	t, err := template.New("").Parse(goTestTemplate)
	assert(err)
	buf := new(bytes.Buffer)
	err = t.Execute(buf, req)
	assert(err)
	return buf.String()
}

var goTestTemplate = `package parse_test

import (
	"testing"

	"github.com/chyroc/dl/internal/parse"
	"github.com/stretchr/testify/assert"
)

func Test_{{ .LowerCamelCaseHost }}(t *testing.T) {
	as := assert.New(t)

	{
		res, err := parse.New{{ .TitleCamelCaseHost }}().Parse("...")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("...", res.Title())
	}
}

`