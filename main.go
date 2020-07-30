package main

import (
	"fmt"
	"github.com/ryer/wpdump/wpdump"
	"github.com/spf13/pflag"
	"os"
)

type appFlags struct {
	help       bool
	url        string
	dir        string
	embed      bool
	merge      bool
	all        bool
	tags       bool
	users      bool
	media      bool
	posts      bool
	pages      bool
	categories bool
	custom     string
}

func parseFlags() *appFlags {
	flags := &appFlags{}

	pflag.BoolVarP(&flags.help, "help", "", false, "show this message")
	pflag.StringVarP(&flags.url, "url", "u", "", "api base url (e.g. http://example.com/wp-json/wp/v2)")
	pflag.StringVarP(&flags.dir, "dir", "d", ".", "save json to this directory")
	pflag.BoolVarP(&flags.embed, "embed", "e", false, "enable embed")
	pflag.BoolVarP(&flags.merge, "merge", "m", false, "merged output")
	pflag.BoolVarP(&flags.all, "all", "a", false, "dump all")
	pflag.BoolVarP(&flags.posts, "posts", "", false, "dump posts")
	pflag.BoolVarP(&flags.categories, "categories", "", false, "dump categories")
	pflag.BoolVarP(&flags.tags, "tags", "", false, "dump tags")
	pflag.BoolVarP(&flags.media, "media", "", false, "dump media")
	pflag.BoolVarP(&flags.pages, "pages", "", false, "dump pages")
	pflag.BoolVarP(&flags.users, "users", "", false, "dump users")
	pflag.StringVarP(&flags.custom, "custom", "", "", "dump custom type")

	pflag.CommandLine.SortFlags = false
	pflag.Parse()

	return flags
}

func main() {
	flags := parseFlags()

	dumpTarget := decideDumpTarget(flags)

	if flags.help || flags.url == "" || len(dumpTarget) == 0 {
		pflag.Usage()
		return
	}

	dumper := buildDumper(flags)
	for _, path := range dumpTarget {
		_, err := dumper.Dump(path)
		if err != nil {
			errorExit(err)
		}
	}
}

func decideDumpTarget(flags *appFlags) []wpdump.Path {
	dumpTarget := make([]wpdump.Path, 0, 7)

	if flags.all || flags.categories {
		dumpTarget = append(dumpTarget, wpdump.Categories)
	}
	if flags.all || flags.pages {
		dumpTarget = append(dumpTarget, wpdump.Pages)
	}
	if flags.all || flags.tags {
		dumpTarget = append(dumpTarget, wpdump.Tags)
	}
	if flags.all || flags.media {
		dumpTarget = append(dumpTarget, wpdump.Media)
	}
	if flags.all || flags.posts {
		dumpTarget = append(dumpTarget, wpdump.Posts)
	}
	if flags.all || flags.users {
		dumpTarget = append(dumpTarget, wpdump.Users)
	}
	if flags.custom != "" {
		dumpTarget = append(dumpTarget, wpdump.Path(flags.custom))
	}

	return dumpTarget
}

func buildDumper(flags *appFlags) wpdump.IDumper {
	var dumper wpdump.IDumper

	if flags.merge {
		dumper = wpdump.NewMergeDumper(wpdump.NewDumper(flags.url, flags.dir, flags.embed), flags.dir)
	} else {
		dumper = wpdump.NewDumper(flags.url, flags.dir, flags.embed)
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
