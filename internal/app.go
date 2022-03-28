package internal

import (
	"fmt"

	"github.com/chyroc/dl/internal/download"
	"github.com/chyroc/dl/internal/identify"
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
	if err := download.Download(args.Dest, resourcer); err != nil {
		return err
	}

	// 5. done
	fmt.Printf("[done] success\n")
	return nil
}

type Argument struct {
	Dest string
	URL  string
}
