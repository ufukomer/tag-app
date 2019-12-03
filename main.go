package main

import (
	"errors"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Usage = "Edit build tags of all go files in specified directory."
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "folder",
			Aliases: []string{"f"},
			Value:   ".",
			Usage:   "target folder",
		},
		&cli.StringFlag{
			Name:    "tags",
			Aliases: []string{"t"},
			Usage:   "tags which are going to be appended",
		},
		&cli.StringFlag{
			Name:    "suffix",
			Aliases: []string{"s"},
			Value:   ".go",
			Usage:   "suffix of target files in target folder",
		},
	}
	app.Before = func(c *cli.Context) error {
		if !c.IsSet("tags") {
			return errors.New("tags flag is not set")
		}
		return nil
	}
	app.Commands = []*cli.Command{
		{
			Name:    "append",
			Aliases: []string{"a"},
			Usage:   "Appends build tag into all files in directory (replaces existing tags)",
			Action: func(c *cli.Context) error {
				return appendTags(c.String("folder"), c.String("tags"), c.String("suffix"))
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
