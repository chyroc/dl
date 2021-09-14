package helper

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
)

var downloadHttpClient = http.Client{}

func Download(url string) (string, error) {
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

	fileSize := resp.ContentLength
	progress := mpb.New(mpb.WithWidth(20))
	bar := progress.AddBar(
		fileSize,
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
	reader := bar.ProxyReader(resp.Body)
	defer reader.Close()

	if _, err := io.Copy(f, reader); err != nil {
		return "", err
	}

	return f.Name(), nil
}

func Download2(url, target string) error {
	file, err := Download(url)
	if err != nil {
		return err
	}
	return os.Rename(file, target)
}
