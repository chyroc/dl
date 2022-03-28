package download

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/chyroc/dl/internal/helper"
	"github.com/chyroc/dl/internal/resource"
	"github.com/mattn/go-runewidth"
)

func Download(dest, prefix string, resourcer resource.Resource) error {
	// 4. download
	switch resourcer := resourcer.(type) {
	case resource.ChapterResource:
		fmt.Printf("%s[chapter] %s\n", prefix, resourcer.Title())
		for idx, v := range resourcer.Chapters() {
			chapterPrefix := fmt.Sprintf("%s[chapter][%d]", prefix, idx+1)
			fmt.Printf("%s %s\n", chapterPrefix, v.Title())
			if err := downloadAtom(filepath.Join(dest, resourcer.Title()), chapterPrefix, v); err != nil {
				return err
			}
		}
	case resource.MP3ChapterResource:
		fmt.Printf("%s[chapter] %s\n", prefix, resourcer.Title())
		for idx, v := range resourcer.Chapters() {
			chapterPrefix := fmt.Sprintf("%s[chapter][%d]", prefix, idx+1)
			fmt.Printf("%s %s\n", chapterPrefix, v.Title())
			if err := downloadAtom(filepath.Join(dest, resourcer.Title()), chapterPrefix, v); err != nil {
				return err
			}
		}
	case resource.Mp3Resource:
		fmt.Printf("%s %s\n", prefix, resourcer.Title())
		if err := downloadAtom(dest, prefix, resourcer); err != nil {
			return err
		}
		if err := resourcer.MP3().UpdateTag(filepath.Join(dest, resourcer.Title())); err != nil {
			return err
		}
	case resource.Resource:
		fmt.Printf("%s %s\n", prefix, resourcer.Title())
		if err := downloadAtom(dest, prefix, resourcer); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupport %T", resourcer)
	}

	return nil
}

func downloadAtom(dest, prefix string, resource resource.Resource) error {
	_ = os.MkdirAll(dest, os.ModePerm)

	tempFile := filepath.Join(dest, resource.Title()+".tmp")
	realFile := filepath.Join(dest, resource.Title())
	f, err := os.OpenFile(tempFile, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	length, reader, err := resource.Reader()
	if err != nil {
		return err
	}

	reader = helper.NewProgressReaderClose(genPrefix(prefix, realFile), length, reader, false)
	defer reader.Close()

	if _, err = io.Copy(f, reader); err != nil {
		return err
	}

	return helper.Rename(tempFile, realFile)
}

func genPrefix(prefix string, filename string) string {
	base := filepath.Base(filename)
	if runewidth.StringWidth(base) > 35 {
		ext := filepath.Ext(base)
		base = base[:len(base)-len(ext)]

		base = runewidth.Truncate(base, 30, "...") + ext
	}

	return fmt.Sprintf("%s[%s] ", prefix, base)
}
