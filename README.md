# wpdump

This is a simple tool that dumps all posts into a json file using wp-json.

## Usage

```
$ wpdump -u 'http://example.com/wp-json/wp/v2' --tags --posts -m
```

```
$ wpdump --help
Usage of wpdump:
  -a, --all          dump all
      --categories   dump categories
  -d, --dir string   save json to this directory (default ".")
  -e, --embed        enable embed
      --help         show this message
      --media        dump media
  -m, --merge        merged output (using jq as an external command)
      --pages        dump pages
      --posts        dump posts
      --tags         dump tags
  -u, --url string   api base url (e.g. http://example.com/wp-json/wp/v2)
      --users        dump users
```

## Installation

```
$ go get github.com/ryer/wpdump
```

## License

MIT

## Author

ryer (@ryer)
