package helper

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/vbauerster/mpb/v7"
)

func DownloadMp3List(mp3List []*MP3) ([]string, error) {
	wg := new(sync.WaitGroup)
	p := mpb.New(mpb.WithWaitGroup(wg), mpb.WithWidth(20))
	num := len(mp3List)

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
			mp3 := mp3List[i]

			f, err := ioutil.TempFile("", "dl-*")
			if err != nil {
				finalErr = err
				return
			}

			req, err := http.NewRequest(http.MethodGet, mp3.DownloadUrl, nil)
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

			err = UpdateMp3Tag(mp3, &mp3.Tag, f.Name())
			if err != nil {
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

func DownloadMp3List2(mp3List []*MP3, dist string) error {
	files, err := DownloadMp3List(mp3List)
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
	return Rename(f.Name(), dist)
}
