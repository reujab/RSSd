package main

import (
	"net"
	"os"
	"os/user"
	"path/filepath"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Usage = "an RSSd client"
	app.HideHelp = true
	app.HideVersion = true
	app.Action = list
	app.Run(os.Args)
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
