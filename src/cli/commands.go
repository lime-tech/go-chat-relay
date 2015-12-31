package cli

import (
	"github.com/codegangsta/cli"
)

var (
	commands = []cli.Command{
		{
			Name:      "start",
			ShortName: "s",
			Usage:     "start a relay",
			Flags:     startFlags,
			Action:    startAction,
		},
	}
)
