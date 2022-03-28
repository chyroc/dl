package internal

import (
	"fmt"
	"path/filepath"

	"github.com/chyroc/dl/internal/download"
	"github.com/chyroc/dl/internal/identify"
	"github.com/chyroc/dl/internal/resource"
)

// Run identify[uri, parser] -> parse[uri, meta] -> download[meta, file]
func DownloadData(args *Argument) error {
	// 2. identify
	parser, err := identify.Identify(args.URL)
	if err != nil {
		return err
	} else if parser == nil {
		return fmt.Errorf("unsupport %q", args.URL)
	}
	fmt.Printf("[meta] %s\n", parser.Kind())

	// 3. parse
	resourcer, err := parser.Parse(args.URL)
	if err != nil {
		return err
	}

	// 4. download
	switch resourcer := resourcer.(type) {
	case resource.ChapterResource:
		fmt.Printf("[chapter] %s\n", resourcer.Title())
		for _, v := range resourcer.Chapters() {
			fmt.Printf("[chapter][download] %s\n", v.Title())
			if err = download.Download(filepath.Join(args.Dest, resourcer.Title()), v); err != nil {
				return err
			}
		}
	case resource.MP3ChapterResource:
		fmt.Printf("[chapter] %s\n", resourcer.Title())
		for _, v := range resourcer.Chapters() {
			fmt.Printf("[chapter][download] %s\n", v.Title())
			if err = download.Download(filepath.Join(args.Dest, resourcer.Title()), v); err != nil {
				return err
			}
		}
	case resource.Mp3Resource:
		fmt.Printf("[download] %s\n", resourcer.Title())
		if err = download.Download(args.Dest, resourcer); err != nil {
			return err
		}
		if err = resourcer.MP3().UpdateTag(filepath.Join(args.Dest, resourcer.Title())); err != nil {
			return err
		}
	case resource.Resource:
		fmt.Printf("[download] %s\n", resourcer.Title())
		if err = download.Download(args.Dest, resourcer); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupport %T", resourcer)
	}

	// 5. done
	fmt.Printf("[done] success\n")
	return nil
}

type Argument struct {
	Dest string
	URL  string
}
