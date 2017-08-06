package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
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
	uris   []*url.URL
	feeds  []*gofeed.Feed
	read   []string
	unread []commands.ListItem
)

func main() {
	usr, err := user.Current()
	die(err)
	file, err := os.Open(filepath.Join(usr.HomeDir, ".config/rssd/feeds"))
	die(err, "no configuration found (~/.config/rssd/feeds)")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		uri, err := url.Parse(scanner.Text())
		die(err, "invalid url: %s", scanner.Text())
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
	defer func() { die(sock.Close()) }()
	go func() {
		for {
			conn, err := sock.Accept()
			die(err)

			reader := bufio.NewReader(conn)
			command, err := reader.ReadByte()
			die(err)

			switch command {
			case commands.List:
				die(json.NewEncoder(conn).Encode(unread))
			case commands.Read:
				var index uint16
				die(binary.Read(reader, binary.BigEndian, &index))
				if int(index) < len(unread) {
					item := unread[index]
					die(json.NewEncoder(conn).Encode(item.URL))
					unread = append(unread[:index], unread[index+1:]...)
					read = append(read, item.GUID)

					// save read GUIDs
					file, err := os.OpenFile(filepath.Join(usr.HomeDir, ".config/rssd/read"), os.O_WRONLY|os.O_CREATE, 0644)
					die(err)
					die(json.NewEncoder(file).Encode(read))
				}
			default:
				die("unknown command")
			}

			die(conn.Close())
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

func update() {
	// update feeds
	var tmpFeeds []*gofeed.Feed
	for _, uri := range uris {
		for i := 0; i < 3; i++ {
			log.Printf("Updating %s...\n", uri)
			parser := gofeed.NewParser()
			feed, err := parser.ParseURL(uri.String())
			if err == nil {
				tmpFeeds = append(tmpFeeds, feed)
				break
			}
		}
	}
	fmt.Println()
	feeds = tmpFeeds

	// update unread articles
	var tmpUnread []commands.ListItem
	for _, feed := range feeds {
	itemLoop:
		for _, item := range feed.Items {
			// check if item has already been read
			for _, guid := range read {
				if guid == item.GUID {
					continue itemLoop
				}
			}

			// item hasn't been read
			tmpUnread = append(tmpUnread, commands.ListItem{
				GUID:  item.GUID,
				Name:  feed.Title,
				Title: item.Title,
				URL:   item.Link,
			})
		}
	}
	unread = tmpUnread
}
