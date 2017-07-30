package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
)

func main() {
	usr, err := user.Current()
	die(err)
	file, err := os.Open(filepath.Join(usr.HomeDir, ".config/feeds"))
	dieMsgIf(err, "no configuration found (~/.config/feeds)")

	var feeds []*url.URL
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		uri, err := url.Parse(scanner.Text())
		dieMsgIf(err, "invalid url: %s", scanner.Text())
		feeds = append(feeds, uri)
	}
	die(scanner.Err())

	fmt.Println(feeds)
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
