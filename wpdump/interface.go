package wpdump

type IDumper interface {
	SetReport(report Report)
	Dump(path Path) ([]string, error)
}

type Report func(path Path, filename string)

type Path string

const (
	PATH_POSTS      = Path("posts")
	PATH_CATEGORIES = Path("categories")
	PATH_TAGS       = Path("tags")
	PATH_MEDIA      = Path("media")
	PATH_PAGES      = Path("pages")
)
