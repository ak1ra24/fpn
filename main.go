package main

import (
	"fmt"
	"os"

	"github.com/ak1ra24/fpn/pkg/cmds"
	"github.com/urfave/cli/v2"
)

const (
	name        = "fpn"
	description = "fake pod network"
	version     = "0.0.1"
)

func main() {
	if err := newApp().Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = name
	app.Version = version
	app.Usage = description
	app.Authors = []*cli.Author{
		{
			Name:  "ak1ra24",
			Email: "ak1ra24net@gmail.com",
		},
	}
	app.Commands = cmds.Commands

	return app
}
