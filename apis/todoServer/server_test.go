package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// setupAPI - used to setup a test server using newMux as main entrypoint
func setupAPI(t *testing.T) (string, func()) {
	t.Helper()

	ts := httptest.NewServer(newMux("")) // create new test server using our mux

	return ts.URL, func() {
		ts.Close() // close/cleanup server function
	}
}

func TestGet(t *testing.T) {
	testCases := []struct {
		name       string
		path       string
		expCode    int
		expItems   int
		expContent string
	}{
		{
			name:       "GetRoot",
			path:       "/",
			expCode:    http.StatusOK,
			expContent: "There's an API here",
		},
		{
			name:    "NotFound",
			path:    "/gibberish",
			expCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// setup test sever
			url, cleanup := setupAPI(t)
			defer cleanup() // cleanup on exit

			var (
				body []byte
				err  error
			)

			r, err := http.Get(url + tc.path)
			if err != nil {
				t.Error(err)
			}
			defer r.Body.Close() // memory leak prevention

			if r.StatusCode != tc.expCode {
				t.Fatalf("Expected %q, got  %q", http.StatusText(tc.expCode), http.StatusText(r.StatusCode))
			}

			switch {
			// if header is plain
			case strings.Contains(r.Header.Get("Content-Type"), "text/plain"):
				// read body and compare bytes
				if body, err = io.ReadAll(r.Body); err != nil {
					t.Error(err)
				}

				// compare body to expected content
				if !strings.Contains(string(body), tc.expContent) {
					t.Errorf("Expected %q, got  %q", tc.expContent, string(body))
				}
			default:
				t.Fatalf("Unsupported Content-Type: %q", r.Header.Get("Content-Type"))
			}

		})
	}
}
