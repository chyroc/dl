package parse

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	"github.com/chyroc/dl/internal/download"
)

func NewMBggeeCom() Parser {
	return &mBggeeCom{}
}

type mBggeeCom struct{}

func (r *mBggeeCom) Kind() string {
	return "m.bggee.com"
}

func (r *mBggeeCom) Parse(uri string) (download.Downloader, error) {
	match := regexp.MustCompile(`https://m.bggee.com/book_(\d+)/`).FindStringSubmatch(uri)
	if len(match) != 2 {
		return nil, fmt.Errorf("匹配不到 book_id")
	}
	bookID := match[1]
	title := ""

	contentURLs := []string{}
	page := 1
	for {
		pageURL := fmt.Sprintf("https://m.bggee.com/index_%s_%d_asc/", bookID, page) // page 1 order asc
		resp, err := http.Get(pageURL)
		if err != nil {
			return nil, err
		}
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		title = regexp.MustCompile(`	var articlename = "(.*?)";`).FindStringSubmatch(string(bs))[1]

		matchs := regexp.MustCompile(`<a class="db" href="(.*?)"`).FindAllStringSubmatch(string(bs), -1)
		for _, v := range matchs {
			if len(v) == 2 {
				contentURLs = append(contentURLs, v[1])
			}
		}

		nextPage := 0
		pageMarch := regexp.MustCompile(`dftval="(\d+)"`).FindStringSubmatch(string(bs))
		if len(pageMarch) == 2 {
			nextPage, _ = strconv.Atoi(pageMarch[1])
		}
		if nextPage == 0 || nextPage != page {
			break
		} else {
			page++
		}
	}

	return download.NewDownloadBggee(bookID, title, contentURLs), nil
}
