package main

import (
	"errors"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Usage = "Edit build tags of all go files in specified directory."
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "folder, f",
			Value: ".",
		},
		cli.StringFlag{
			Name: "tags, t",
		},
		cli.StringFlag{
			Name:  "suffix, s",
			Value: ".go",
		},
	}
	app.Before = func(c *cli.Context) error {
		if !c.GlobalIsSet("tags") {
			return errors.New("tags flag is not set")
		}
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:  "append, a",
			Usage: "Appends build tag into all files in directory (replaces existing tags)",
			Action: func(c *cli.Context) error {
				return appendTags(c.GlobalString("folder"), c.GlobalString("tags"), c.GlobalString("suffix"))
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
