package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/ryer/wpdump/wpdump"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		expectedExit   int
		expectedOutput string
	}{
		{
			name:           "version",
			args:           []string{"--version"},
			expectedExit:   0,
			expectedOutput: "wpdump version",
		},
		{
			name:           "help",
			args:           []string{"--help"},
			expectedExit:   1,
			expectedOutput: "Usage of",
		},
		{
			name:           "no url",
			args:           []string{"--posts"},
			expectedExit:   1,
			expectedOutput: "Usage of",
		},
		{
			name:           "no target",
			args:           []string{"--url", "http://example.com"},
			expectedExit:   1,
			expectedOutput: "Usage of",
		},
		{
			name: "success run",
			args: []string{
				"--url", "DUMMY_URL", // Placeholder, will be replaced by mock server URL
				"--posts",
				"--dir",
			},
			expectedExit:   0,
			expectedOutput: "Dump end",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock server for success case
			var server *httptest.Server
			if tc.name == "success run" {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Header().Set("X-Wp-Totalpages", "1")
					_, err := w.Write([]byte(`[{"id":1}]`))
					assert.NoError(t, err)
				}))
				defer server.Close()

				// Replace placeholder URL and add temp dir
				tempDir := t.TempDir()
				for i, arg := range tc.args {
					if arg == "DUMMY_URL" {
						tc.args[i] = server.URL
					}
				}
				tc.args = append(tc.args, tempDir)
			}

			// Backup and defer restore
			originalArgs := os.Args
			originalStdout := os.Stdout
			originalStderr := os.Stderr
			defer func() {
				os.Args = originalArgs
				os.Stdout = originalStdout
				os.Stderr = originalStderr
				pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
			}()

			// Redirect stdout/stderr
			r, w, _ := os.Pipe()
			os.Stdout = w
			os.Stderr = w

			// Set args and run
			os.Args = append([]string{"wpdump"}, tc.args...)
			pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
			exitCode := run()

			w.Close()
			var buf bytes.Buffer
			_, err := io.Copy(&buf, r)
			require.NoError(t, err)

			if exitCode != tc.expectedExit {
				t.Errorf("Expected exit code %d, but got %d", tc.expectedExit, exitCode)
			}
			if !strings.Contains(buf.String(), tc.expectedOutput) {
				t.Errorf("Expected output to contain '%s', but got '%s'", tc.expectedOutput, buf.String())
			}
		})
	}
}

func TestDecideDumpTarget(t *testing.T) {
	flags := &appFlags{}

	flags.categories = true
	a := decideDumpTarget(flags)

	if len(a) != 1 {
		t.Fatalf("did not select a 1 path (%v)", len(a))
	}

	if a[0] != wpdump.Categories {
		t.Fatalf("did not select categories (%v)", a[0])
	}

	flags.all = true
	a = decideDumpTarget(flags)

	if len(a) != 6 {
		t.Fatalf("did not select 6 paths (%v)", len(a))
	}
}

func TestBuildDumper(t *testing.T) {
	flags := &appFlags{}

	flags.merge = false
	a, _ := buildDumper(flags)

	if reflect.TypeOf(a).String() != "*wpdump.WPDumper" {
		t.Fatalf("is not WPDump (%v)", reflect.TypeOf(a))
	}

	flags.merge = true
	a, _ = buildDumper(flags)

	if reflect.TypeOf(a).String() != "*wpdump.WPMergeDumper" {
		t.Fatalf("is not MergeDumper (%v)", reflect.TypeOf(a))
	}
}
