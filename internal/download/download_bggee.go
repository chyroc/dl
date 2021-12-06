package download

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func NewDownloadBggee(booID, title string, contentURLs []string) Downloader {
	return &downloadBggee{
		bookID:      booID,
		title:       title,
		contentURLs: contentURLs,
	}
}

type downloadBggee struct {
	bookID      string
	title       string
	contentURLs []string
}

func (r *downloadBggee) Title() string {
	return fmt.Sprintf("%s_%s.txt", r.title, r.bookID)
}

func (r *downloadBggee) TargetFile() string {
	return ""
}

func (r *downloadBggee) Download() error {
	textList := []string{}
	for _, contentURL := range r.contentURLs {
		resp, err := http.Get(contentURL)
		if err != nil {
			return err
		}
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		d, _ := goquery.NewDocumentFromReader(bytes.NewReader(bs))
		textList = append(textList, d.Find(".acontent").Text())
	}
	return ioutil.WriteFile(r.Title(), []byte(strings.Join(textList, "\n")), 0644)
}

func (r *downloadBggee) MultiDownload() []Downloader {
	return nil
}
