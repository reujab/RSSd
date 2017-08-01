# RSSd
RSSd is a daemon that reads RSS feeds.

## Getting started
### Installation
```sh
go get github.com/reujab/RSSd/...
```

This will install two binaries:
* `RSSd`, the daemon
* `rss`, the client

### Adding feeds
Add all the RSS feed URLs to `~/.config/rssd/feeds` separated by newlines.

Example:
```
http://eslint.org/feed.xml
https://blog.golang.org/feed.atom
https://electron.atom.io/feed.xml
https://fedoramagazine.org/feed/
https://github.com/showcases.atom
https://www.youtube.com/feeds/videos.xml?channel_id=UCsLiV4WJfkTEHH0b9PmRklw
https://xkcd.com/atom.xml
```

### Using
Before you can use the CLI, you must start the daemon with `RSSd`.

#### Listing unread items
```
$ rss
1   GitHub Showcases  Great for new contributors
2   GitHub Showcases  Government apps
3   GitHub Showcases  Open Source Integrations
4   GitHub Showcases  Web accessibility
5   GitHub Showcases  Social Impact
6   GitHub Showcases  Programming languages
7   GitHub Showcases  GitHub Browser Extensions
8   GitHub Showcases  Game Engines
9   GitHub Showcases  Software Defined Radio
10  GitHub Showcases  Text editors
11  GitHub Showcases  JavaScript game engines
12  GitHub Showcases  Clean code linters
13  GitHub Showcases  DevOps tools
14  GitHub Showcases  Virtual Reality
15  GitHub Showcases  Tools for Open Source
```

#### Opening an item
```
$ rss 6
```

#### Marking all items as read
```
$ rss read-all
```
