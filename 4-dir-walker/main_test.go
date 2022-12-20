package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// createTempDir
func createTempDir(t *testing.T, files map[string]int) (dirname string, cleanup func()) {
	t.Helper() // this is a helper function that can be called by other tests inside same package

	// create a new temporary directory
	tempDir, err := os.MkdirTemp("", "walktest")
	if err != nil {
		t.Fatal(err)
	}

	// loop over each request file and create n number of files
	for k, n := range files {
		for j := 1; j <= n; j++ {
			fname := fmt.Sprintf("file%d%s", j, k)
			fpath := filepath.Join(tempDir, fname) // join directory + filename
			if err := os.WriteFile(fpath, []byte("dummy"), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}

	return tempDir, func() { os.RemoveAll(tempDir) }
}

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		root     string
		cfg      config
		expected string
	}{
		{
			name: "NoFilter",
			root: "testdata",
			cfg: config{
				ext:  "",
				size: 0,
				list: true,
			},
			expected: "testdata/dir.log\ntestdata/dir2/script.sh\n",
		},
		{
			name: "FilterExtensionMatch",
			root: "testdata",
			cfg: config{
				ext:  ".log",
				size: 0,
				list: true,
			},
			expected: "testdata/dir.log\n",
		},
		{
			name: "FilterExtensionSizeMatch",
			root: "testdata",
			cfg: config{
				ext:  ".log",
				size: 10,
				list: true,
			},
			expected: "testdata/dir.log\n",
		},
		{
			name: "FilterExtensionSizeNoMatch",
			root: "testdata",
			cfg: config{
				ext:  ".log",
				size: 20,
				list: true,
			},
			expected: "",
		},
		{
			name: "FilterExtensionNoMatch",
			root: "testdata",
			cfg: config{
				ext:  ".gz",
				size: 0,
				list: true,
			},
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer

			if err := run(tc.root, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			res := buffer.String()

			if res != tc.expected {
				t.Errorf("got %q, want %q", res, tc.expected)
			}

		})
	}

}

func TestRunDeleteExtension(t *testing.T) {
	testCases := []struct {
		name        string
		cfg         config
		extNoDelete string
		nDelete     int
		nNoDelete   int
		expected    string
	}{
		{
			name: "DeleteExtensionNoMatch",
			cfg: config{
				ext: ".log",
				del: true,
			},
			extNoDelete: ".gz",
			nDelete:     0,
			nNoDelete:   10,
			expected:    "",
		},
		{
			name: "DeleteExtensionMatch",
			cfg: config{
				ext: ".log",
				del: true,
			},
			extNoDelete: "",
			nDelete:     10,
			nNoDelete:   0,
			expected:    "",
		},
		{
			name: "DeleteExtensionMixed",
			cfg: config{
				ext: ".log",
				del: true,
			},
			extNoDelete: ".gz",
			nDelete:     5,
			nNoDelete:   5,
			expected:    "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer

			// create temporary directory and files
			tempDir, cleanup := createTempDir(t, map[string]int{
				tc.cfg.ext:     tc.nDelete,
				tc.extNoDelete: tc.nNoDelete,
			})
			// fmt.Println("temp directory created at: ", tempDir)
			defer cleanup() //remove tempDir on exit

			// run deletion on tempDir
			if err := run(tempDir, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			// compare results
			got := buffer.String()
			want := tc.expected
			if got != want {
				t.Errorf("got %q, want %q\n", got, want)
			}

			filesLeft, err := os.ReadDir(tempDir)
			if err != nil {
				t.Error(err)
			}

			if len(filesLeft) != tc.nNoDelete {
				t.Errorf("Expected %d files left, got %d instead\n", tc.nNoDelete, len(filesLeft))
			}
		})
	}
}
