package wpdump

import "time"

type IDumper interface {
	SetReporter(reporter Reporter)
	Dump(path Path) ([]string, error)
}

type Reporter interface {
	Start()
	End(elapsed time.Duration)
	Success(path Path, filename string)
	Error(err error)
	Warn(message string)
}

type Path string

const (
	Posts      = Path("posts")
	Categories = Path("categories")
	Tags       = Path("tags")
	Media      = Path("media")
	Pages      = Path("pages")
	Users      = Path("users")
)
