package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/chyroc/dl/pkgs"
	"github.com/chyroc/dl/pkgs/helper"
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
		Action: func(c *cli.Context) error {
			args, err := parseArgument(c)
			if err != nil {
				return err
			}
			return pkgs.DownloadData(args)
		},
		Commands: []*cli.Command{
			{
				Name: "example",
				Action: func(c *cli.Context) error {
					pkgs.RunExample(c)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func parseArgument(c *cli.Context) (*pkgs.Argument, error) {
	dest, err := helper.ResolveDirOrCurrent(strings.TrimSpace(c.String("dest")))
	if err != nil {
		return nil, err
	}

	uri := c.Args().Get(0)
	if uri == "" {
		return nil, fmt.Errorf("must set uri to download")
	}

	return &pkgs.Argument{
		Dest: dest,
		URL:  uri,
	}, nil
}
