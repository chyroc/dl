package internal

import (
	"fmt"

	"github.com/chyroc/dl/internal/identify"
)

func RunExample() {
	for _, url := range identify.ExampleURLs {
		if url == "" {
			continue
		}
		err := DownloadData(&Argument{Dest: "/tmp/", URL: url})
		if err != nil {
			fmt.Printf("[fail] download %q failed: %v\n", url, err)
		} else {
			fmt.Printf("[succ] download %q success\n", url)
		}
		// fmt.Println()
	}
}
