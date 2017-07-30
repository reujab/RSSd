package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"syscall"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/reujab/RSSd/commands"
)

var (
	uris  []*url.URL
	feeds []*gofeed.Feed
	read  []string
)

func main() {
	usr, err := user.Current()
	die(err)
	file, err := os.Open(filepath.Join(usr.HomeDir, ".config/rssd/feeds"))
	dieMsgIf(err, "no configuration found (~/.config/rssd/feeds)")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		uri, err := url.Parse(scanner.Text())
		dieMsgIf(err, "invalid url: %s", scanner.Text())
		uris = append(uris, uri)
	}
	die(scanner.Err())

	file, err = os.Open(filepath.Join(usr.HomeDir, ".config/rssd/read"))
	if err == nil {
		die(json.NewDecoder(file).Decode(&read))
	}

	go func() {
		for {
			update()
			time.Sleep(time.Minute * 10)
		}
	}()

	sock, err := net.Listen("unix", filepath.Join(usr.HomeDir, ".local/share/rssd.sock"))
	die(err)
	defer func() { sock.Close() }()
	go func() {
		for {
			conn, err := sock.Accept()
			die(err)
			defer func() { die(conn.Close()) }()

			reader := bufio.NewReader(conn)
			command, err := reader.ReadByte()
			die(err)

			switch command {
			case commands.List:
				fmt.Println("requested list")
			default:
				die("unknown command")
			}
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt
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

func dieMsg(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func dieMsgIf(err error, format string, args ...interface{}) {
	if err != nil {
		dieMsg(format, args...)
	}
}

func update() {
	var tmpFeeds []*gofeed.Feed
	for _, uri := range uris {
		parser := gofeed.NewParser()
		feed, err := parser.ParseURL(uri.String())
		die(err)
		tmpFeeds = append(tmpFeeds, feed)
	}
	feeds = tmpFeeds
}
