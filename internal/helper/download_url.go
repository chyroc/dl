package helper

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
)

var downloadHttpClient = http.Client{}

func Download(url string, chapter bool) (string, error) {
	f, err := ioutil.TempFile("", "dl-*")
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := downloadHttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	reader := newProgressReaderClose(resp.ContentLength, resp.Body, chapter)
	defer reader.Close()

	if _, err := io.Copy(f, reader); err != nil {
		return "", err
	}

	return f.Name(), nil
}

func Download2(url, target string, chapter bool) error {
	file, err := Download(url, chapter)
	if err != nil {
		return err
	}
	return Rename(file, target)
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
