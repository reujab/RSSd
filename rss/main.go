package main

import (
	"net"
	"os"
	"os/user"
	"path/filepath"
	"strconv"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Usage = "an RSSd client"
	app.HideHelp = true
	app.HideVersion = true
	app.Commands = []cli.Command{
		{
			Name:   "read-all",
			Usage:  "marks all items as read",
			Action: readAll,
		},
	}
	app.Action = func(ctx *cli.Context) {
		switch len(ctx.Args()) {
		case 0:
			list()
		case 1:
			index, err := strconv.ParseUint(ctx.Args()[0], 10, 16)
			die(err, "invalid index")
			read(uint16(index) - 1)
		default:
			die("too many arguments")
		}
	}
	die(app.Run(os.Args))
}

func die(args ...interface{}) {
	switch len(args) {
	case 0:
		os.Exit(1)
	case 1:
		if args[0] != nil {
			panic(args[0])
		}
	case 2:
		if args[0] != nil {
			panic(args[1])
		}
	}
}

func connect() net.Conn {
	usr, err := user.Current()
	die(err)
	conn, err := net.Dial("unix", filepath.Join(usr.HomeDir, ".local/share/rssd.sock"))
	die(err)
	return conn
}
