package download

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mattn/go-runewidth"

	"github.com/chyroc/dl/pkgs/helper"
	"github.com/chyroc/dl/pkgs/resource"
)

func Download(dest, prefix string, resourcer resource.Resourcer) error {
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
	case resource.Resourcer:
		fmt.Printf("%s %s\n", prefix, resourcer.Title())
		if err := downloadAtom(dest, prefix, resourcer); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupport %T", resourcer)
	}

	return nil
}

func downloadAtom(dest, prefix string, resourceIns resource.Resourcer) error {
	_ = os.MkdirAll(dest, os.ModePerm)
	realFile := filepath.Join(dest, resourceIns.Title())

	temFile, err := os.CreateTemp("", "dl-output-temp-*.mp4")
	if err != nil {
		return err
	}
	var length int64
	var lengthGen func() int64
	var reader io.Reader
	if r, ok := resourceIns.(resource.Resourcer2); ok {
		lengthGen, reader, err = r.Reader2()
	} else {
		length, reader, err = resourceIns.Reader()
	}
	if err != nil {
		return err
	}

	readCloser := helper.NewProgressReaderClose(genPrefix(prefix, realFile), length, lengthGen, reader, false)
	defer readCloser.Close()

	if _, err = io.Copy(temFile, readCloser); err != nil {
		return err
	}

	temFilepath := temFile.Name()
	if triggerIns, ok := resourceIns.(resource.ResourcerAfterTrigger); ok {
		temFilepath, err = triggerIns.Trigger(temFilepath)
		if err != nil {
			return err
		}
	}

	return helper.Rename(temFilepath, realFile)
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
