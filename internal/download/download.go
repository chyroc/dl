package download

import (
	"io"
	"os"
	"path/filepath"

	"github.com/chyroc/dl/internal/helper"
	"github.com/chyroc/dl/internal/resource"
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

	reader = helper.NewProgressReaderClose(length, reader, false)
	defer reader.Close()

	if _, err = io.Copy(f, reader); err != nil {
		return err
	}

	return helper.Rename(tempFile, filepath.Join(dest, resource.Title()))
}
