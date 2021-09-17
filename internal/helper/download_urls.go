package helper

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
)

func DownloadURLs(urls []string) ([]string, error) {
	wg := new(sync.WaitGroup)
	p := mpb.New(mpb.WithWaitGroup(wg), mpb.WithWidth(20))
	num := len(urls)

	// var bars = make([]*mpb.Bar, numBars)
	files := make([]string, num)
	var finalErr error
	for i := 0; i < num; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			if finalErr != nil {
				return
			}
			url := urls[i]

			f, err := ioutil.TempFile("", "dl-*")
			if err != nil {
				finalErr = err
				return
			}

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				finalErr = err
				return
			}

			resp, err := downloadHttpClient.Do(req)
			if err != nil {
				finalErr = err
				return
			}
			defer resp.Body.Close()

			reader := newProgressReaderCloseV7(p, resp.ContentLength, resp.Body)
			defer reader.Close()

			if _, err = io.Copy(f, reader); err != nil {
				finalErr = err
				return
			}

			files[i] = f.Name()
			// bars = append(bars, b)
		}(i)
	}

	p.Wait()

	return files, nil
}

func DownloadURLs2(urls []string, dist string) error {
	files, err := DownloadURLs(urls)
	if err != nil {
		return err
	}

	f, err := ioutil.TempFile("", "dl-*")
	if err != nil {
		return err
	}

	for _, file := range files {
		temp, err := os.OpenFile(file, os.O_RDONLY, 0o666)
		if err != nil {
			return err
		}
		if _, err = f.ReadFrom(temp); err != nil {
			return err
		}
	}
	return os.Rename(f.Name(), dist)
}

func newProgressReaderCloseV7(p *mpb.Progress, length int64, body io.Reader) io.ReadCloser {
	bar := p.AddBar(
		length,
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
	return bar.ProxyReader(body)
}
