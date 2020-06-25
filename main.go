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
		users      = flag.Bool("users", false, "Dump users")
		all        = flag.Bool("all", false, "Dump all")
		embed      = flag.Bool("embed", false, "Enable embed")
		merge      = flag.Bool("merge", false, "Merged output (using jq as an external command)")
	)
	flag.Parse()

	pathList := make([]wpdump.Path, 0, 6)
	if *all || *categories {
		pathList = append(pathList, wpdump.PATH_CATEGORIES)
	}
	if *all || *pages {
		pathList = append(pathList, wpdump.PATH_PAGES)
	}
	if *all || *tags {
		pathList = append(pathList, wpdump.PATH_TAGS)
	}
	if *all || *media {
		pathList = append(pathList, wpdump.PATH_MEDIA)
	}
	if *all || *posts {
		pathList = append(pathList, wpdump.PATH_POSTS)
	}
	if *all || *users {
		pathList = append(pathList, wpdump.PATH_USERS)
	}

	dumper := buildDumper(*url, *dir, *embed, *merge)
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

func buildDumper(baseUrl string, outputDir string, embed bool, merge bool) wpdump.IDumper {
	var dumper wpdump.IDumper

	if merge {
		dumper = wpdump.NewMergeDumper(wpdump.NewDumper(baseUrl, outputDir, embed), outputDir)
	} else {
		dumper = wpdump.NewDumper(baseUrl, outputDir, embed)
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
