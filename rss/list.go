package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/reujab/RSSd/commands"
)

func list() {
	conn := connect()
	conn.Write([]byte{commands.List})

	var items []commands.ListItem
	die(json.NewDecoder(conn).Decode(&items))

	lengths := make([]int, 2)
	for i, item := range items {
		if len(item.Name) > lengths[1] {
			lengths[1] = len(item.Name)
		}

		lengths[0] = len(strconv.Itoa(i + 1))
	}

	for i, item := range items {
		var line string

		line += strconv.Itoa(i + 1)
		line += strings.Repeat(" ", lengths[0]-len(strconv.Itoa(i+1))) + "  "

		line += item.Name
		line += strings.Repeat(" ", lengths[1]-len(item.Name)) + "  "

		line += item.Title

		fmt.Println(line)
	}
}
