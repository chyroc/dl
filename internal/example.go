package internal

import (
	"fmt"

	"github.com/chyroc/dl/internal/identify"
	"github.com/urfave/cli/v2"
)

func RunExample(c *cli.Context) {
	url := ""
	if c.Args().Len() > 0 {
		url = c.Args().Get(0)
	}
	if url != "" {
		downloadExample(url)
	} else {
		for _, url := range identify.ExampleURLs {
			if url == "" {
				continue
			}
			downloadExample(url)
		}
	}
}

func downloadExample(url string) {
	err := DownloadData(&Argument{Dest: "/tmp/", URL: url})
	if err != nil {
		fmt.Printf("[fail] download %q failed: %v\n", url, err)
	} else {
		fmt.Printf("[succ] download %q success\n", url)
	}
}
