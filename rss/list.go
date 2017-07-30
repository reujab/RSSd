package main

import (
	"github.com/reujab/RSSd/commands"
	"github.com/urfave/cli"
)

func list(ctx *cli.Context) {
	conn := connect()
	conn.Write([]byte{commands.List})
}
