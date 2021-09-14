package internal

import (
	"fmt"

	"github.com/chyroc/dl/internal/identify"
	"github.com/urfave/cli/v2"
)

// identify[uri, parser] -> parse[uri, meta] -> download[meta, file]
func Run(c *cli.Context) error {
	uri := c.Args().Get(0)
	if uri == "" {
		return fmt.Errorf("must set uri to download")
	}

	parser, err := identify.Identify(uri)
	if err != nil {
		return err
	} else if parser == nil {
		return fmt.Errorf("unsupport %q", uri)
	}

	downloader, err := parser.Parse(uri)
	if err != nil {
		return err
	}

	err = downloader.Download()
	if err != nil {
		return err
	}

	return nil
}
