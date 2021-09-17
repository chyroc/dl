package main

import (
	"io/ioutil"
	"sort"
	"strings"
)

func main() {
	dir := "./internal/parse/"
	fs, err := ioutil.ReadDir(dir)
	assert(err)
	urls := []string{}
	for _, v := range fs {
		if strings.HasSuffix(v.Name(), "_test.go") {
			continue
		}
		bs, err := ioutil.ReadFile(dir + "" + v.Name())
		assert(err)
		list := strings.Split(string(bs), "\n")
		for idx, vv := range list {
			vv = strings.TrimSpace(vv)
			if strings.HasSuffix(vv, "Kind() string {") {
				vvv := strings.TrimSpace(list[idx+1])
				vvv = vvv[len("return "):]
				vvv = vvv[1 : len(vvv)-1]
				urls = append(urls, vvv)
			}
		}
	}

	sort.Strings(urls)

	bs, err := ioutil.ReadFile("./README.md")
	assert(err)
	s := string(bs)
	s1 := strings.Split(s, "## Support Website")[0] + "## Support Website"
	s2 := "## Install" + strings.Split(s, "## Install")[1]

	readme := s1 + "\n\n"
	for _, v := range urls {
		readme += "- " + v + "\n"
	}
	readme += "\n" + s2
	assert(ioutil.WriteFile("./README.md", []byte(readme), 0o666))
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}
