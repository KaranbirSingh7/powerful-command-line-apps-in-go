package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/karanbirsingh7/pscan/scan"
)

func setup(t *testing.T, hosts []string, initList bool) (string, func()) {
	// t.Helper() // signifies this is a helper func (also silence the errors)
	tf, err := os.CreateTemp("", "pScan")
	if err != nil {
		t.Fatal(err) // fail, tests cannot continue
	}

	// initialize a temp list if needed
	if initList {
		hl := &scan.HostsList{}
		for _, h := range hosts {
			hl.Add(h)
		}

		if err := hl.Save(tf.Name()); err != nil {
			t.Fatal(err)
		}
	}

	return tf.Name(), func() {
		os.Remove(tf.Name())
	}
}

func TestHostActions(t *testing.T) {
	// define hosts for actions test
	hosts := []string{
		"host1",
		"host2",
		"host3",
	}

	testCases := []struct {
		name           string
		args           []string
		expectedOut    string
		initList       bool
		actionFunction func(io.Writer, string, []string) error
	}{
		{
			name:           "AddAction",
			args:           hosts,
			expectedOut:    "Added host: host1\nAdded host: host2\nAdded host: host3\n",
			initList:       false,
			actionFunction: addActions,
		},
		{
			name:           "DeleteAction",
			args:           []string{"host1", "host2"},
			expectedOut:    "Deleted host: host1\nDeleted host: host2\n",
			initList:       true,
			actionFunction: deleteAction,
		},
		{
			name:           "ListAction",
			expectedOut:    "host1\nhost2\nhost3\n",
			args:           hosts,
			initList:       true,
			actionFunction: listAction,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// setup basic
			tempFile, cleanup := setup(t, tc.args, tc.initList)
			defer cleanup()

			// capture output
			var out bytes.Buffer

			// TEST assert/check values
			if err := tc.actionFunction(&out, tempFile, tc.args); err != nil {
				t.Fatalf("Expected no error, got %q\n", err)
			}

			if out.String() != tc.expectedOut {
				t.Errorf("Expected %q, got %q", tc.expectedOut, out.String())
			}
		})
	}
}

// simulate what a user would do
func TestIntegration(t *testing.T) {
	hosts := []string{
		"host1",
		"host2",
		"host3",
	}

	tf, cleanup := setup(t, hosts, false) // creates an empty file /shrug
	defer cleanup()                       //delete file on exit

	delHost := "host2"

	hostsEnd := []string{
		"host1",
		"host3",
	}

	// capture out as buffer
	var out bytes.Buffer

	// define expected output for all actions
	expectedOut := ""
	for _, v := range hosts {
		expectedOut += fmt.Sprintf("Added host: %s\n", v) // add 3 hosts
	}
	expectedOut += strings.Join(hosts, "\n") // list all hosts
	expectedOut += fmt.Sprintln()
	expectedOut += fmt.Sprintf("Deleted host: %s\n", delHost) // delete host2
	expectedOut += strings.Join(hostsEnd, "\n")               // list only 2 hosts since we removed 1
	expectedOut += fmt.Sprintln()

	// TEST start
	// ADD 3 hosts
	if err := addActions(&out, tf, hosts); err != nil {
		t.Fatalf("Expected no error, got %q", err)
	}

	// LIST
	if err := listAction(&out, tf, nil); err != nil {
		t.Fatalf("Expected no error, got %q", err)
	}

	// DELETE
	if err := deleteAction(&out, tf, []string{delHost}); err != nil {
		t.Fatalf("Expected no error, got %q", err)
	}

	// LIST
	if err := listAction(&out, tf, nil); err != nil {
		t.Fatalf("Expected no error, got %q", err)
	}

	// VERIFY
	if out.String() != expectedOut {
		t.Errorf("Expected output %q, got %q", expectedOut, out.String())
	}
}
