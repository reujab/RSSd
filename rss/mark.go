package main

import (
	"encoding/binary"
	"encoding/json"
	"sort"
	"strconv"

	"github.com/reujab/RSSd/commands"
	"github.com/urfave/cli"
)

func mark(ctx *cli.Context) {
	var items []uint16
	for _, arg := range ctx.Args() {
		item, err := strconv.ParseInt(arg, 10, 16)
		die(err)
		items = append(items, uint16(item))
	}

	// greatest to least
	sort.Slice(items, func(i, j int) bool {
		return items[i] > items[j]
	})

	for _, item := range items {
		conn := connect()
		conn.Write([]byte{commands.Read})
		die(binary.Write(conn, binary.BigEndian, item-1))
		die(json.NewDecoder(conn).Decode(new(string)))
	}
}
