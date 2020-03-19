package cmds

import "github.com/urfave/cli/v2"

var Commands = []*cli.Command{
	commandConf,
	commandDown,
	commandPrint,
	commandReConf,
	commandReUp,
	commandUp,
	commandUpConf,
	commandTest,
}

var commandPrint = &cli.Command{
	Name:   "print",
	Usage:  "print",
	Action: CmdPrint,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Specify the Config file.",
			Value:   "spec.yaml",
		},
	},
}

var commandUp = &cli.Command{
	Name:   "up",
	Usage:  "up",
	Action: CmdUp,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Specify the Config file.",
			Value:   "spec.yaml",
		},
	},
}

var commandConf = &cli.Command{
	Name:   "conf",
	Usage:  "conf",
	Action: CmdConf,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Specify the Config file.",
			Value:   "spec.yaml",
		},
	},
}

var commandDown = &cli.Command{
	Name:   "down",
	Usage:  "down",
	Action: CmdDown,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Specify the Config file.",
			Value:   "spec.yaml",
		},
	},
}

var commandUpConf = &cli.Command{
	Name:   "upconf",
	Usage:  "upconf",
	Action: CmdUpConf,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Specify the Config file.",
			Value:   "spec.yaml",
		},
	},
}

var commandReUp = &cli.Command{
	Name:   "reup",
	Usage:  "reup",
	Action: CmdReUp,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Specify the Config file.",
			Value:   "spec.yaml",
		},
	},
}

var commandReConf = &cli.Command{
	Name:   "reconf",
	Usage:  "reconf",
	Action: CmdReConf,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Specify the Config file.",
			Value:   "spec.yaml",
		},
	},
}

var commandTest = &cli.Command{
	Name:   "test",
	Usage:  "test",
	Action: CmdTest,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Specify the Config file.",
			Value:   "spec.yaml",
		},
	},
}
