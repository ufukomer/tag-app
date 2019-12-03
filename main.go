package main

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name: "append",
			Flags: []cli.Flag{
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
			},
			Usage: "append build tag into all files in directory, replaces existing tags.",
			Action: func(c *cli.Context) error {
				if !c.IsSet("tags") {
					return errors.New("tags flag is not set")
				}

				return appendTags(c.String("folder"), c.String("tags"), c.String("suffix"))
			},
		},
	}
}
