# wpdump

This is a simple tool that dumps all posts into a json file using wp-json.

## Usage

```
$ wpdump -u 'http://example.com/wp-json/wp/v2' --tags --posts --custom books -m
```

```
$ wpdump
Usage of wpdump:
      --help                 show this message
      --version              show version
  -u, --url string           api base url (e.g. http://example.com/wp-json/wp/v2)
  -d, --dir string           save json to this directory (default ".")
  -p, --parallel int         parallel download (default 1)
  -v, --verbose              verbose output
  -e, --embed                enable embed
  -m, --merge                merged output
  -a, --all                  dump all
      --posts                dump posts
      --categories           dump categories
      --tags                 dump tags
      --media                dump media
      --pages                dump pages
      --users                dump users
      --custom stringArray   dump custom type (support multiple flags)
```

## Installation

```
$ go install github.com/ryer/wpdump@latest
```

## License

MIT

## Author

ryer (@ryer)
