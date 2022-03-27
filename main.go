package main

import (
	"log"
	"os"

	"github.com/chyroc/dl/internal"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "dl",
		Usage: "Download Chinese Website Video, Audio, Image, Document, etc.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "dest",
				Aliases: []string{"d"},
				Usage:   "set destination directory(default: current directory)",
			},
		},
		Action: internal.Run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
