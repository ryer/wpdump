package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ryer/wpdump/wpdump"
	"github.com/spf13/pflag"
)

var (
	Name     = "wpdump"
	Version  = "1.0.0"
	Revision = "latest"
)

type appFlags struct {
	help       bool
	version    bool
	url        string
	dir        string
	verbose    bool
	parallel   int
	embed      bool
	merge      bool
	all        bool
	tags       bool
	users      bool
	media      bool
	posts      bool
	pages      bool
	categories bool
	custom     []string
}

func parseFlags() *appFlags {
	flags := &appFlags{}

	pflag.BoolVarP(&flags.help, "help", "", false, "show this message")
	pflag.BoolVarP(&flags.version, "version", "", false, "show version")
	pflag.StringVarP(&flags.url, "url", "u", "", "api base url (e.g. http://example.com/wp-json/wp/v2)")
	pflag.StringVarP(&flags.dir, "dir", "d", ".", "save json to this directory")
	pflag.IntVarP(&flags.parallel, "parallel", "p", 1, "parallel download")
	pflag.BoolVarP(&flags.verbose, "verbose", "v", false, "verbose output")
	pflag.BoolVarP(&flags.embed, "embed", "e", false, "enable embed")
	pflag.BoolVarP(&flags.merge, "merge", "m", false, "merged output")
	pflag.BoolVarP(&flags.all, "all", "a", false, "dump all")
	pflag.BoolVarP(&flags.posts, "posts", "", false, "dump posts")
	pflag.BoolVarP(&flags.categories, "categories", "", false, "dump categories")
	pflag.BoolVarP(&flags.tags, "tags", "", false, "dump tags")
	pflag.BoolVarP(&flags.media, "media", "", false, "dump media")
	pflag.BoolVarP(&flags.pages, "pages", "", false, "dump pages")
	pflag.BoolVarP(&flags.users, "users", "", false, "dump users")
	pflag.StringArrayVarP(&flags.custom, "custom", "", []string{}, "dump custom type (support multiple flags)")
	pflag.CommandLine.SortFlags = false
	pflag.Parse()

	return flags
}

func main() {
	flags := parseFlags()

	dumpTarget := decideDumpTarget(flags)

	if flags.version {
		fmt.Printf("%v version %v (%v)", Name, Version, Revision)

		return
	}

	if flags.help || flags.url == "" || len(dumpTarget) == 0 {
		pflag.Usage()

		return
	}

	dumper, reporter := buildDumper(flags)
	reporter.Start()
	start := time.Now()

	for _, path := range dumpTarget {
		_, err := dumper.Dump(path)
		if err != nil {
			// Error is already reported by the reporter
			os.Exit(1)
		}
	}
	reporter.End(time.Since(start))
}

func decideDumpTarget(flags *appFlags) []wpdump.Path {
	targets := map[wpdump.Path]bool{
		wpdump.Categories: flags.categories,
		wpdump.Pages:      flags.pages,
		wpdump.Tags:       flags.tags,
		wpdump.Media:      flags.media,
		wpdump.Posts:      flags.posts,
		wpdump.Users:      flags.users,
	}

	dumpTarget := make([]wpdump.Path, 0, len(targets))
	for path, enabled := range targets {
		if flags.all || enabled {
			dumpTarget = append(dumpTarget, path)
		}
	}

	for _, path := range flags.custom {
		dumpTarget = append(dumpTarget, wpdump.Path(path))
	}

	return dumpTarget
}

type ConsoleReporter struct {
	verbose bool
}

func (r *ConsoleReporter) Start() {
	fmt.Println("Dump start")
}

func (r *ConsoleReporter) End(elapsed time.Duration) {
	if r.verbose {
		fmt.Printf("Dump end (elapsed: %v)\n", elapsed)
	} else {
		fmt.Println("Dump end")
	}
}

func (r *ConsoleReporter) Success(path wpdump.Path, filename string) {
	if r.verbose {
		fmt.Printf("Dump progress (%v): %v\n", path, filename)
	}
}

func (r *ConsoleReporter) Error(err error) {
	// Errors are always displayed
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
}

func (r *ConsoleReporter) Warn(message string) {
	if r.verbose {
		fmt.Fprintf(os.Stderr, "Warn: %s\n", message)
	}
}

func buildDumper(flags *appFlags) (wpdump.IDumper, wpdump.Reporter) {
	var dumper wpdump.IDumper
	reporter := &ConsoleReporter{verbose: flags.verbose}

	dumperImpl := wpdump.NewDumper(flags.url, flags.dir, flags.embed)

	if flags.parallel > 1 {
		dumper = wpdump.NewParallelDumper(dumperImpl, flags.parallel)
	} else {
		dumper = dumperImpl
	}

	if flags.merge {
		dumper = wpdump.NewMergeDumper(dumper, flags.dir)
	}

	dumper.SetReporter(reporter)

	return dumper, reporter
}
