package internal

import (
	"fmt"
	"strings"

	"github.com/chyroc/dl/internal/identify"
	"github.com/urfave/cli/v2"
)

func RunExample(c *cli.Context) {
	url := ""
	if c.Args().Len() > 0 {
		url = c.Args().Get(0)
	}
	if url != "" {
		if strings.HasPrefix(url, "http") {
			downloadExample(url)
		} else {
			for _, internalURL := range identify.ExampleURLs {
				if strings.Contains(internalURL, url) {
					downloadExample(internalURL)
				}
			}
		}
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
	defer func() {
		e := recover()
		if e != nil {
			fmt.Printf("[panic] download %q failed: %v\n", url, e)
		}
	}()
	err := DownloadData(&Argument{Dest: "/tmp/", URL: url})
	if err != nil {
		fmt.Printf("[fail] download %q failed: %v\n", url, err)
	} else {
		fmt.Printf("[succ] download %q success\n", url)
	}
}
