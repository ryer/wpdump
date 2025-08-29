package wpdump

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParallelDumper_Dump(t *testing.T) {
	const totalPages = 5
	const parallel = 2

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("page")
		if pageStr == "" {
			pageStr = "1"
		}
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			http.Error(w, "Invalid page number", http.StatusBadRequest)

			return
		}

		w.Header().Set("X-WP-TotalPages", strconv.Itoa(totalPages))
		fmt.Fprintf(w, `[{"id":%d}]`, page)
	}))
	defer server.Close()

	tempDir, err := os.MkdirTemp("", "wpdump_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	dumper := NewDumper(server.URL, tempDir, false)
	assert.NoError(t, err)

	pd := NewParallelDumper(dumper, parallel)

	path := Path("posts")
	files, err := pd.Dump(path)
	assert.NoError(t, err)
	assert.Len(t, files, totalPages)

	for i := 1; i <= totalPages; i++ {
		expectedFile := filepath.Join(tempDir, fmt.Sprintf("%s%04d.json", path, i))
		assert.FileExists(t, expectedFile)
	}
}

func TestParallelDumper_Dump_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer server.Close()

	tempDir, err := os.MkdirTemp("", "wpdump_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	dumper := NewDumper(server.URL, tempDir, false)
	assert.NoError(t, err)

	pd := NewParallelDumper(dumper, 2)

	path := Path("posts")
	_, err = pd.Dump(path)
	assert.Error(t, err)
}
