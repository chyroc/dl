package download

import (
	"github.com/chyroc/dl/internal/helper"
)

func NewDownloadURLList(title string, targetFile string, urls []string) Downloader {
	return &downloadURLList{
		title:      title,
		targetFile: targetFile,
		urls:       urls,
	}
}

type downloadURLList struct {
	title      string
	targetFile string
	urls       []string
}

func (r *downloadURLList) Title() string {
	return r.title
}

func (r *downloadURLList) TargetFile() string {
	return r.targetFile
}

func (r *downloadURLList) Download() error {
	return helper.DownloadURLs2(r.urls, r.targetFile)
}

func (r *downloadURLList) MultiDownload() []Downloader {
	return nil
}
