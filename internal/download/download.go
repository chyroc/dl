package download

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/chyroc/dl/internal/helper"
	"github.com/chyroc/dl/internal/resource"
)

func Download(dest string, resourcer resource.Resource) error {
	// 4. download
	switch resourcer := resourcer.(type) {
	case resource.ChapterResource:
		fmt.Printf("[chapter] %s\n", resourcer.Title())
		for _, v := range resourcer.Chapters() {
			fmt.Printf("[chapter][download] %s\n", v.Title())
			if err := downloadAtom(filepath.Join(dest, resourcer.Title()), v); err != nil {
				return err
			}
		}
	case resource.MP3ChapterResource:
		fmt.Printf("[chapter] %s\n", resourcer.Title())
		for _, v := range resourcer.Chapters() {
			fmt.Printf("[chapter][download] %s\n", v.Title())
			if err := downloadAtom(filepath.Join(dest, resourcer.Title()), v); err != nil {
				return err
			}
		}
	case resource.Mp3Resource:
		fmt.Printf("[download] %s\n", resourcer.Title())
		if err := downloadAtom(dest, resourcer); err != nil {
			return err
		}
		if err := resourcer.MP3().UpdateTag(filepath.Join(dest, resourcer.Title())); err != nil {
			return err
		}
	case resource.Resource:
		fmt.Printf("[download] %s\n", resourcer.Title())
		if err := downloadAtom(dest, resourcer); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupport %T", resourcer)
	}
}

func downloadAtom(dest string, resource resource.Resource) error {
	_ = os.MkdirAll(dest, os.ModePerm)

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
