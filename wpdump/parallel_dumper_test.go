package wpdump

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

		w.Header().Set("X-Wp-Totalpages", strconv.Itoa(totalPages))
		fmt.Fprintf(w, `[{"id":%d}]`, page)
	}))
	defer server.Close()

	tempDir := t.TempDir()

	dumper := NewDumper(server.URL, tempDir, false)

	pd := NewParallelDumper(dumper, parallel)

	path := Path("posts")
	files, err := pd.Dump(path)
	require.NoError(t, err)
	assert.Len(t, files, totalPages)

	for i := 1; i <= totalPages; i++ {
		expectedFile := filepath.Join(tempDir, fmt.Sprintf("%s%04d.json", path, i))
		assert.FileExists(t, expectedFile)
	}
}

func TestParallelDumper_Dump_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer server.Close()

	tempDir := t.TempDir()

	dumper := NewDumper(server.URL, tempDir, false)

	pd := NewParallelDumper(dumper, 2)

	path := Path("posts")
	_, err := pd.Dump(path)
	require.Error(t, err)
}

func TestParallelDumper_Dump_LessPagesThanParallel(t *testing.T) {
	const totalPages = 1
	const parallel = 2

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Wp-Totalpages", strconv.Itoa(totalPages))
		fmt.Fprintf(w, `[{"id":1}]`)
	}))
	defer server.Close()

	tempDir := t.TempDir()

	dumper := NewDumper(server.URL, tempDir, false)
	pd := NewParallelDumper(dumper, parallel)

	path := Path("posts")
	files, err := pd.Dump(path)
	require.NoError(t, err)
	assert.Len(t, files, totalPages)

	expectedFile := filepath.Join(tempDir, fmt.Sprintf("%s%04d.json", path, 1))
	assert.FileExists(t, expectedFile)
}

func TestParallelDumper_Dump_PartialError(t *testing.T) {
	const totalPages = 5
	const parallel = 2
	const errorPage = 3

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("page")
		page, _ := strconv.Atoi(pageStr)

		if page == errorPage {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)

			return
		}

		w.Header().Set("X-Wp-Totalpages", strconv.Itoa(totalPages))
		fmt.Fprintf(w, `[{"id":%d}]`, page)
	}))
	defer server.Close()

	tempDir := t.TempDir()

	dumper := NewDumper(server.URL, tempDir, false)
	pd := NewParallelDumper(dumper, parallel)

	path := Path("posts")
	_, err := pd.Dump(path)
	require.Error(t, err)
}
