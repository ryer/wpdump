# wpdump

This is a simple tool that dumps all posts into a json file using wp-json.

## Usage

```
$ wpdump -u 'http://example.com/wp-json/wp/v2' --tags --posts -m
```

```
$ wpdump --help
Usage of wpdump:
      --help            show this message
  -u, --url string      api base url (e.g. http://example.com/wp-json/wp/v2)
  -d, --dir string      save json to this directory (default ".")
  -e, --embed           enable embed
  -m, --merge           merged output
  -a, --all             dump all
      --posts           dump posts
      --categories      dump categories
      --tags            dump tags
      --media           dump media
      --pages           dump pages
      --users           dump users
      --custom string   dump custom type
```

## Installation

```
$ go get github.com/ryer/wpdump
```

## License

MIT

## Author

ryer (@ryer)
