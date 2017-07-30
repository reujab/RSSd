package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/k0kubun/pp"
	"github.com/mmcdole/gofeed"
)

var feeds []*url.URL

func main() {
	usr, err := user.Current()
	die(err)
	file, err := os.Open(filepath.Join(usr.HomeDir, ".config/rssd/feeds"))
	dieMsgIf(err, "no configuration found (~/.config/rssd/feeds)")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		uri, err := url.Parse(scanner.Text())
		dieMsgIf(err, "invalid url: %s", scanner.Text())
		feeds = append(feeds, uri)
	}
	die(scanner.Err())

	for {
		update()
	}
}

func die(err error) {
	if err != nil {
		panic(err)
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
	for _, uri := range feeds {
		parser := gofeed.NewParser()
		feed, err := parser.ParseURL(uri.String())
		die(err)
		pp.Println(feed)
	}

	time.Sleep(time.Minute * 10)
}
