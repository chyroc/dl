package download

import (
	"io"
	"os"
	"path/filepath"

	"github.com/chyroc/dl/internal/helper"
	"github.com/chyroc/dl/internal/resource"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
)

type Downloader interface {
	Title() string
	TargetFile() string
	Download() error
	MultiDownload() []Downloader
}

func Download(dest string, resource resource.Resource) error {
	tempFile := filepath.Join(dest, resource.Title()+".tmp")
	f, err := os.OpenFile(tempFile, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	length, reader, err := resource.Reader()
	if err != nil {
		return err
	}

	reader = newProgressReaderClose(length, reader, false)
	defer reader.Close()

	if _, err = io.Copy(f, reader); err != nil {
		return err
	}

	return helper.Rename(tempFile, filepath.Join(dest, resource.Title()))
}

func newProgressReaderClose(length int64, body io.Reader, chapter bool) io.ReadCloser {
	progress := mpb.New(mpb.WithWidth(20))
	barOptions := []mpb.BarOption{}
	if true {
		barOptions = append(barOptions, mpb.BarRemoveOnComplete())
	}
	barOptions = append(barOptions,
		// 进度条前的修饰
		mpb.PrependDecorators(
			decor.Name("[download] "),
			decor.CountersKibiByte("% .2f / % .2f"), // 已下载数量
			decor.Percentage(decor.WCSyncSpace),     // 进度百分比
		),
		// 进度条后的修饰
		mpb.AppendDecorators(
			decor.EwmaSpeed(decor.UnitKiB, "% .2f", 60),
		),
	)
	bar := progress.AddBar(
		length,
		barOptions...,
	)
	return bar.ProxyReader(body)
}
