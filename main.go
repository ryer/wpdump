package main

import (
	"fmt"
	"github.com/ryer/wpdump/wpdump"
	"github.com/spf13/pflag"
	"os"
)

func main() {
	var (
		url        = pflag.StringP("url", "u", "", "api base url (e.g. http://example.com/wp-json/wp/v2)")
		dir        = pflag.StringP("dir", "d", ".", "save json to this directory")
		posts      = pflag.BoolP("posts", "", false, "dump posts")
		categories = pflag.BoolP("categories", "", false, "dump categories")
		tags       = pflag.BoolP("tags", "", false, "dump tags")
		media      = pflag.BoolP("media", "", false, "dump media")
		pages      = pflag.BoolP("pages", "", false, "dump pages")
		users      = pflag.BoolP("users", "", false, "dump users")
		all        = pflag.BoolP("all", "a", false, "dump all")
		embed      = pflag.BoolP("embed", "e", false, "enable embed")
		merge      = pflag.BoolP("merge", "m", false, "merged output (using jq as an external command)")
	)
	pflag.Parse()

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
		pflag.Usage()
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
