package download

import (
	"fmt"
	"os"

	"github.com/chyroc/goexec"
	"github.com/google/uuid"
)

func NewDownloadM3U8(title string, targetFile string, url string) Downloader {
	return &downloadM3U8{
		title:      title,
		targetFile: targetFile,
		url:        url,
	}
}

type downloadM3U8 struct {
	title      string
	targetFile string
	url        string
}

func (r *downloadM3U8) Title() string {
	return r.title
}

func (r *downloadM3U8) TargetFile() string {
	return r.targetFile
}

func (r *downloadM3U8) Download() error {
	f := fmt.Sprintf("%s%s.mp4", os.TempDir(), uuid.New().String())
	_, _, err := goexec.New("ffmpeg", "-i", r.url, "-c", "copy", "-f", "mp4", f).Run()
	if err != nil {
		return err
	}

	return os.Rename(f, r.targetFile)
}
