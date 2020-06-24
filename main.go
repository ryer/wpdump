package main

import (
	"flag"
	"fmt"
	"github.com/ryer/wpdump/wpdump"
	"os"
)

func main() {
	var (
		url        = flag.String("url", "", "API Base URL (e.g. http://example.com/wp-json/wp/v2)")
		dir        = flag.String("dir", ".", "Save json to this directory")
		posts      = flag.Bool("posts", false, "Dump posts")
		categories = flag.Bool("categories", false, "Dump categories")
		tags       = flag.Bool("tags", false, "Dump tags")
		media      = flag.Bool("media", false, "Dump media")
		pages      = flag.Bool("pages", false, "Dump pages")
		merge      = flag.Bool("merge", false, "Merged output (using jq as an external command)")
	)
	flag.Parse()

	pathList := make([]wpdump.Path, 0, 5)
	if *categories {
		pathList = append(pathList, wpdump.PATH_CATEGORIES)
	}
	if *pages {
		pathList = append(pathList, wpdump.PATH_PAGES)
	}
	if *tags {
		pathList = append(pathList, wpdump.PATH_TAGS)
	}
	if *media {
		pathList = append(pathList, wpdump.PATH_MEDIA)
	}
	if *posts {
		pathList = append(pathList, wpdump.PATH_POSTS)
	}

	dumper := buildDumper(*url, *dir, *merge)
	for _, path := range pathList {
		_, err := dumper.Dump(path)
		if err != nil {
			errorExit(err)
		}
	}

	if len(pathList) == 0 {
		flag.Usage()
	}
}

func buildDumper(baseUrl string, outputDir string, merge bool) wpdump.IDumper {
	var dumper wpdump.IDumper

	if merge {
		dumper = wpdump.NewMergeDumper(baseUrl, outputDir)
	} else {
		dumper = wpdump.NewDumper(baseUrl, outputDir)
	}

	dumper.SetReport(func(path wpdump.Path, filename string) {
		fmt.Println(fmt.Sprintf("Dumped(%v): %v", path, filename))
	})

	return dumper
}

func errorExit(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}
