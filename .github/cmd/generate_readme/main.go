package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/chyroc/dl/.github/cmd/cmd_helper"
)

func main() {
	dir := "./pkgs/parse/"
	fs, err := ioutil.ReadDir(dir)
	cmd_helper.Assert(err)

	urls := []string{}
	for _, v := range fs {
		if v.IsDir() {
			continue
		}
		if strings.HasSuffix(v.Name(), "_test.go") {
			continue
		}
		bs, err := ioutil.ReadFile(dir + "" + v.Name())
		cmd_helper.Assert(err)
		list := strings.Split(string(bs), "\n")
		for idx, vv := range list {
			vv = strings.TrimSpace(vv)
			if strings.HasSuffix(vv, "Kind() string {") {
				vvv := strings.TrimSpace(list[idx+1])
				vvv = vvv[len("return "):]
				vvv = vvv[1 : len(vvv)-1]

				for _, vvvv := range strings.Split(vvv, ",") {
					urls = append(urls, vvvv)
				}
			}
		}
	}

	sort.Strings(urls)

	bs, err := ioutil.ReadFile("./README.md")
	cmd_helper.Assert(err)

	s := string(bs)
	s1 := strings.Split(s, "## Support Website")[0] + "## Support Website"
	s2 := "## Install" + strings.Split(s, "## Install")[1]

	readme := fmt.Sprintf("%s\n\nsupport %d websites\n\n", s1, len(urls))
	for _, v := range urls {
		readme += fmt.Sprintf("- [%s](%s)\n", v, v)
	}
	readme += "\n" + s2
	cmd_helper.Assert(ioutil.WriteFile("./README.md", []byte(readme), 0o666))
}
