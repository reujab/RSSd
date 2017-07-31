package main

import (
	"encoding/binary"
	"encoding/json"
	"os/exec"
	"runtime"

	"github.com/reujab/RSSd/commands"
)

func read(index uint16) {
	conn := connect()
	conn.Write([]byte{commands.Read})

	die(binary.Write(conn, binary.BigEndian, index))

	var uri string
	err := json.NewDecoder(conn).Decode(&uri)
	if err != nil {
		die("index out of bounds")
	}

	switch runtime.GOOS {
	case "darwin":
		die(exec.Command("open", uri).Run())
	default:
		die(exec.Command("xdg-open", uri).Run())
	}
}
