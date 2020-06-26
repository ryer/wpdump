package wpdump

import (
	"github.com/jarcoal/httpmock"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestDump(t *testing.T) {
	mockJSON := `[{"id":3,"name":"ryer","url":"","description":"","link":"https://example.com/author/ryer/","slug":"ryer","avatar_urls":{"24":"https://example.com/com.jpg","48":"https://example.com/bar.jpg","96":"https://example.com/foo.jpg"},"meta":[],"_links":{"self":[{"href":"https://example.com/wp-json/wp/v2/users/3"}],"collection":[{"href":"https://example.com/wp-json/wp/v2/users"}]}}]` + "\n"
	dumper := NewMockDumper(mockJSON)

	files, err := dumper.Dump(Users)
	if err != nil {
		t.Fatalf("an error occurred (%v)", err)
	}

	if files[0] != "/tmp/users0001.json" {
		t.Fatalf("file name mismatch (%v)", files[0])
	}

	data, err := ioutil.ReadFile(files[0])
	if err != nil {
		t.Fatalf("an error occurred (%v)", err)
	}

	if string(data) != mockJSON {
		t.Fatalf("data mismatch (%v)", string(data))
	}
}

func NewMockDumper(mockJSON string) *WPDumper {
	dumper := NewDumper("https://example.com/wp-json/wp/v2", "/tmp", false)

	httpmock.ActivateNonDefault(dumper.client.GetClient())
	mockHeader := http.Header{}
	mockHeader.Add("X-WP-TotalPages", "1")
	responder := httpmock.ResponderFromResponse(&http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Body:          httpmock.NewRespBodyFromString(mockJSON),
		Header:        mockHeader,
		ContentLength: int64(len(mockJSON)),
	})
	fakeURL := "https://example.com/wp-json/wp/v2/users"
	httpmock.RegisterResponder("GET", fakeURL, responder)

	return dumper
}
