package internal

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/chyroc/dl/internal/download"
	"github.com/chyroc/dl/internal/helper"
	"github.com/chyroc/dl/internal/identify"
	"github.com/chyroc/dl/internal/resource"
	"github.com/urfave/cli/v2"
)

// Run identify[uri, parser] -> parse[uri, meta] -> download[meta, file]
func Run(c *cli.Context) error {
	// 1. arguments
	args, err := parseArgument(c)
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
	case resource.Resource:
		fmt.Printf("[download] %s\n", resourcer.Title())
		if err = download.Download(args.Dest, resourcer); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupport %T", resourcer)
	}

	// 5. done
	fmt.Println("[done] success")
	return nil
}

type Argument struct {
	Dest string
	URL  string
}

func parseArgument(c *cli.Context) (*Argument, error) {
	dest, err := helper.ResolveDirOrCurrent(strings.TrimSpace(c.String("dest")))
	if err != nil {
		return nil, err
	}

	uri := c.Args().Get(0)
	if uri == "" {
		return nil, fmt.Errorf("must set uri to download")
	}

	return &Argument{
		Dest: dest,
		URL:  uri,
	}, nil
}
