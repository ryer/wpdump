package wpdump

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestMergeDump(t *testing.T) {
	mockJSON := `[{"id":3,"name":"ryer","url":"","description":"","link":"https://example.com/author/ryer/","slug":"ryer","avatar_urls":{"24":"https://example.com/com.jpg","48":"https://example.com/bar.jpg","96":"https://example.com/foo.jpg"},"meta":[],"_links":{"self":[{"href":"https://example.com/wp-json/wp/v2/users/3"}],"collection":[{"href":"https://example.com/wp-json/wp/v2/users"}]}}]`
	dumper := NewMergeDumper(NewMockDumper(mockJSON), "./testdata")

	files, err := dumper.Dump(Users)
	if err != nil {
		t.Fatalf("an error occurred (%v)", err)
	}

	if files[0] != "./testdata/users.json" {
		t.Fatalf("file name mismatch (%v)", files[0])
	}

	data, err := ioutil.ReadFile(files[0])
	if err != nil {
		t.Fatalf("an error occurred (%v)", err)
	}

	innerJSON := strings.Trim(mockJSON, "[]")
	expected := "[" + innerJSON + ",\n" + innerJSON + "]"
	if string(data) != expected {
		t.Fatalf("data mismatch (%v)", string(data))
	}
}
