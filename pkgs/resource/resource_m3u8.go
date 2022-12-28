package resource

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/oopsguy/m3u8/parse"
	"github.com/oopsguy/m3u8/tool"

	"github.com/chyroc/dl/pkgs/helper"
)

type m3u8Resource struct {
	concurrency int
	title       string
	url         string
}

func NewM3U8(concurrency int, title string, url string) (Resourcer, error) {
	return &m3u8Resource{
		concurrency: concurrency,
		title:       title,
		url:         url,
	}, nil
}

func (r *m3u8Resource) Title() string {
	return r.title
}

func (r *m3u8Resource) Reader() (int64, io.ReadCloser, error) {
	panic("not impl")
}

func (r *m3u8Resource) Trigger(file string) (string, error) {
	return helper.M3u8ToMp4(file)
}

func (r *m3u8Resource) Reader2() (func() int64, io.ReadCloser, error) {
	result, err := parse.FromURL(r.url)
	if err != nil {
		return nil, nil, fmt.Errorf("parse m3u8: %s fail: %s", r.url, err.Error())
	}

	segLen := len(result.M3u8.Segments)
	res := &reader{
		concurrency: r.concurrency,
		finish:      0,
		segLen:      segLen,
		result:      result,
	}
	return func() int64 {
		return res.totalLength
	}, res, nil
}

type reader struct {
	concurrency   int
	totalLength   int64
	lock          sync.Mutex
	finish        int32
	segLen        int
	result        *parse.Result
	readers       *sync.Map
	currentReader int
	once          sync.Once
}

func (mr *reader) Read(p []byte) (n int, err error) {
	mr.once.Do(func() {
		mr.readers = new(sync.Map)
		go mr.fetchData()
	})

	for {
		if mr.currentReader >= mr.segLen {
			return 0, io.EOF
		}
		reader_, ok := mr.readers.Load(mr.currentReader)
		if !ok {
			time.Sleep(time.Second)
			continue
		}
		reader := reader_.(io.Reader)

		n, err = reader.Read(p)
		if err == io.EOF {
			mr.readers.Delete(mr.currentReader)
			mr.currentReader++
		}
		if n > 0 || err != io.EOF {
			if err == io.EOF && mr.currentReader < mr.segLen {
				err = nil
			}
			return
		}
	}
	return 0, io.EOF
}

func (r *reader) Close() error {
	return nil
}

func (r *reader) fetchData() {
	var wg sync.WaitGroup
	limitChan := make(chan struct{}, r.concurrency)
	for tsIndex := 0; tsIndex < r.segLen; tsIndex++ {
		tsIndex := tsIndex
		wg.Add(1)
		go func(tsIndex int) {
			defer wg.Done()

			for {
				length, oneReader, err := r.download(tsIndex)
				if err != nil {
					continue
				}
				r.readers.Store(tsIndex, oneReader)
				r.totalLength += length
				break
			}
			<-limitChan
		}(tsIndex)
		limitChan <- struct{}{}
	}
	wg.Wait()
}

func (d *reader) download(segIndex int) (int64, io.Reader, error) {
	tsUrl := d.tsURL(segIndex)
	b, err := getURL(tsUrl)
	if err != nil {
		return 0, nil, fmt.Errorf("request %s, %s", tsUrl, err.Error())
	}
	defer b.Close()

	body, err := io.ReadAll(b)
	if err != nil {
		return 0, nil, fmt.Errorf("read body: %s, %s", tsUrl, err.Error())
	}
	sf := d.result.M3u8.Segments[segIndex]
	if sf == nil {
		return 0, nil, fmt.Errorf("invalid segment index: %d", segIndex)
	}
	key, ok := d.result.Keys[sf.KeyIndex]
	if ok && key != "" {
		body, err = tool.AES128Decrypt(body, []byte(key), []byte(d.result.M3u8.Keys[sf.KeyIndex].IV))
		if err != nil {
			return 0, nil, fmt.Errorf("decryt: %s, %s", tsUrl, err.Error())
		}
	}
	// https://en.wikipedia.org/wiki/MPEG_transport_stream
	// Some TS files do not start with SyncByte 0x47, they can not be played after merging,
	// Need to remove the body before the SyncByte 0x47(71).
	syncByte := uint8(71) // 0x47
	bLen := len(body)
	for j := 0; j < bLen; j++ {
		if body[j] == syncByte {
			body = body[j:]
			break
		}
	}

	atomic.AddInt32(&d.finish, 1)
	return int64(len(body)), bytes.NewReader(body), nil
}

func (d *reader) tsURL(segIndex int) string {
	seg := d.result.M3u8.Segments[segIndex]
	return tool.ResolveURL(d.result.URL, seg.URI)
}

func getURL(url string) (io.ReadCloser, error) {
	resp, err := downloadHttpClient.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http error: status code %d", resp.StatusCode)
	}
	return resp.Body, nil
}
