package main

import (
	"io/ioutil"
	"sort"
	"strings"
)

func main() {
	dir := "internal/parse/"
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
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

	// package identify
	//
	// import (
	//	"github.com/chyroc/dl/internal/parse"
	// )
	//
	// func init() {
	//	register(parse.NewA36krCom())
	// }
	content := []string{
		"package identify",
		"",
		"import (",
		"	\"github.com/chyroc/dl/internal/parse\"",
		")",
		"",
		"func init() {",
	}
	for _, v := range funcNames {
		content = append(content, "	register(parse."+v+"())")
	}
	content = append(content, "}")

	err = ioutil.WriteFile("internal/identify/register.go", []byte(strings.Join(content, "\n")), 0644)
	if err != nil {
		panic(err)
	}
}

func handle(path string) string {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	// func NewYQqCom() Parser {
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
