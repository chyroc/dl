package download

import (
	"github.com/chyroc/dl/internal/helper"
)

func NewDownloadURL(title string, targetFile string, pkgs []*Specification) Downloader {
	return &downloadURL{
		title:      title,
		targetFile: targetFile,
		specs:      pkgs,
	}
}

type downloadURL struct {
	title      string
	targetFile string
	specs      []*Specification
}

func (r *downloadURL) Title() string {
	return r.title
}

func (r *downloadURL) TargetFile() string {
	return r.targetFile
}

func (r *downloadURL) Download() error {
	pkg := SpecificationList(r.specs).GetMax()

	return helper.Download2(pkg.URL, r.TargetFile())
}
