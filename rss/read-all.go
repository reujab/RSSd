package main

import (
	"encoding/binary"
	"encoding/json"

	"github.com/reujab/RSSd/commands"
	"github.com/urfave/cli"
)

func readAll(ctx *cli.Context) {
	for {
		conn := connect()
		conn.Write([]byte{commands.Read})

		die(binary.Write(conn, binary.BigEndian, uint16(0)))

		var uri string
		err := json.NewDecoder(conn).Decode(&uri)
		if err != nil {
			break
		}
	}
}
