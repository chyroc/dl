package main

import (
	"io/ioutil"
	"sort"
	"strings"

	"github.com/chyroc/dl/.github/cmd/cmd_helper"
)

func main() {
	dir := "pkgs/parse/"
	fs, err := ioutil.ReadDir(dir)
	cmd_helper.Assert(err)

	funcNames := []string{}
	for _, f := range fs {
		if f.IsDir() {
			continue
		}
		if f.Name()[len(f.Name())-3:] != ".go" {
			continue
		}
		funcName := handle(dir + f.Name())
		if funcName != "" {
			funcNames = append(funcNames, funcName)
		}
	}
	sort.Strings(funcNames)

	content := []string{
		"package identify",
		"",
		"import (",
		"	\"github.com/chyroc/dl/pkgs/parse\"",
		")",
		"",
		"func init() {",
	}
	for _, v := range funcNames {
		content = append(content, "	register(parse."+v+"())")
	}
	content = append(content, "}")

	err = ioutil.WriteFile("pkgs/identify/register.go", []byte(strings.Join(content, "\n")), 0o644)
	cmd_helper.Assert(err)
}

func handle(path string) string {
	bs, err := ioutil.ReadFile(path)
	cmd_helper.Assert(err)

	s := string(bs)
	for _, v := range strings.Split(s, "\n") {
		v = strings.TrimSpace(v)
		if strings.HasPrefix(v, "func New") && strings.HasSuffix(v, "() Parser {") {
			funcName := v[len("func ") : len(v)-len("() Parser {")]
			return funcName
		}
	}
	return ""
}
