package wpdump

type IDumper interface {
	SetReport(report Report)
	Dump(path Path) ([]string, error)
}

type Report func(path Path, filename string)

type Path string

const (
	Posts      = Path("posts")
	Categories = Path("categories")
	Tags       = Path("tags")
	Media      = Path("media")
	Pages      = Path("pages")
	Users      = Path("users")
)
