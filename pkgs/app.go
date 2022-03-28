package pkgs

import (
	"fmt"
	"net/url"

	"github.com/chyroc/dl/pkgs/download"
	"github.com/chyroc/dl/pkgs/identify"
)

// Run identify[uri, parser] -> parse[uri, meta] -> download[meta, file]
func DownloadData(args *Argument) error {
	// 1. host
	URL, err := url.Parse(args.URL)
	if err != nil {
		return err
	}

	// 2. identify
	parser, err := identify.Identify(args.URL)
	if err != nil {
		return err
	} else if parser == nil {
		return fmt.Errorf("unsupport %q", args.URL)
	}
	parser.Kind()
	fmt.Printf("[%s] resource kind=%s\n", URL.Host, parser.Kind())

	// 3. parse
	resourcer, err := parser.Parse(args.URL)
	if err != nil {
		return err
	}

	// 4. download
	if err := download.Download(args.Dest, "["+URL.Host+"]", resourcer); err != nil {
		return err
	}

	// 5. done
	fmt.Printf("[%s] download %q success\n", URL.Host, args.URL)
	return nil
}

type Argument struct {
	Dest string
	URL  string
}
