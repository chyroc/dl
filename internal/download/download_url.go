package download

import (
	"github.com/chyroc/dl/internal/helper"
)

func NewDownloadURL(title string, targetFile string, chapter bool, pkgs []*Specification) Downloader {
	return &downloadURL{
		title:      title,
		targetFile: targetFile,
		specs:      pkgs,
		chapter:    chapter,
	}
}

type downloadURL struct {
	title      string
	targetFile string
	specs      []*Specification
	chapter    bool
}

func (r *downloadURL) Title() string {
	return r.title
}

func (r *downloadURL) TargetFile() string {
	return r.targetFile
}

func (r *downloadURL) Download() error {
	pkg := SpecificationList(r.specs).GetMax()

	return helper.Download2(pkg.URL, r.TargetFile(), r.chapter)
}

func (r *downloadURL) MultiDownload() []Downloader {
	return nil
}
