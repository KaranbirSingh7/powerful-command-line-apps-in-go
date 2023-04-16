package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/karanbirsingh7/pclaig/todo"
)

// setupAPI - used to setup a test server using newMux as main entrypoint
func setupAPI(t *testing.T) (string, func()) {
	t.Helper()

	tempTodoFile, err := os.CreateTemp("", "todotest")
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(newMux(tempTodoFile.Name())) // create new test server using our mux

	// seed todofile with items
	for i := 1; i < 3; i++ {
		var body bytes.Buffer
		taskName := fmt.Sprintf("Task number %d.", i)
		item := struct {
			Task string `json:"task"`
		}{
			Task: taskName,
		}

		// encode json
		if err := json.NewEncoder(&body).Encode(item); err != nil {
			t.Fatal(err)
		}
		// make POST call to add item
		r, err := http.Post(ts.URL+"/todo", "application/json", &body)
		if err != nil {
			t.Fatal(err)
		}
		// check response
		if r.StatusCode != http.StatusCreated {
			t.Fatalf("Failed to add initial items: Status: %d", r.StatusCode)
		}
	}

	return ts.URL, func() {
		ts.Close()                     // close/cleanup server function
		os.Remove(tempTodoFile.Name()) //delete seeded todolist file
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
		{
			name:       "GetAll",
			path:       "/todo",
			expCode:    http.StatusOK,
			expItems:   2,
			expContent: "Task number 1.",
		},
		{
			name:       "GetOne",
			path:       "/todo/1",
			expCode:    http.StatusOK,
			expItems:   1,
			expContent: "Task number 1.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// setup test sever
			url, cleanup := setupAPI(t)
			defer cleanup() // cleanup on exit

			var (
				resp struct {
					Results      todo.List `json:"results"`
					Date         int64     `json:"date"`
					TotalResults int       `json:"total_results"`
				}
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
			case r.Header.Get("Content-Type") == "application/json":
				if err = json.NewDecoder(r.Body).Decode(&resp); err != nil {
					t.Error(err)
				}
				if resp.TotalResults != tc.expItems {
					t.Errorf("Expected %d items, got %d.", tc.expItems, resp.TotalResults)
				}
				if resp.Results[0].Task != tc.expContent {
					t.Errorf("Expected %q, got %q.", tc.expContent,
						resp.Results[0].Task)
				}
			case strings.Contains(r.Header.Get("Content-Type"), "text/plain"):
				if body, err = io.ReadAll(r.Body); err != nil {
					t.Error(err)
				}

				if !strings.Contains(string(body), tc.expContent) {
					t.Errorf("Expected %q, got %q.", tc.expContent,
						string(body))
				}
			default:
				t.Fatalf("Unsupported Content-Type: %q", r.Header.Get("Content-Type"))
			}

		})
	}
}
