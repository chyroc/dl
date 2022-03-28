package helper

import (
	"io"

	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
)

type ProgressReader struct {
	io.ReadCloser

	progress *mpb.Progress
	bar      *mpb.Bar
}

func (r *ProgressReader) SetTotal(n int64, complete bool) {
	r.bar.SetTotal(n, complete)
}

func NewProgressReaderClose(length int64, body io.Reader, chapter bool) *ProgressReader {
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
	return &ProgressReader{
		ReadCloser: bar.ProxyReader(body),
		progress:   progress,
		bar:        bar,
	}
}
