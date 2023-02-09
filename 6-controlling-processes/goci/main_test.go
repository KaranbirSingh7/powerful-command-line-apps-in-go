package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func setupGit(t *testing.T, proj string) func() {
	t.Helper() // this is a helper function that can be called inside other test(s), this helps with errors line numbers.

	// check if git cli is present
	gitExec, err := exec.LookPath("git")
	if err != nil {
		// fatal error since we cannot continue if git is not installed
		t.Fatal(err)
	}

	// create a temp dir
	tempDir, err := os.MkdirTemp("", "gocitest")
	if err != nil {
		t.Fatal(err)
	}

	projPath, err := filepath.Abs(proj)
	if err != nil {
		t.Fatal(err)
	}

	// remote URI for git repo (this is interesting approach), we are pointing this to a temporary directory
	remoteURI := fmt.Sprintf("file://%s", tempDir)

	// execute step of commands
	var gitCmdList = []struct {
		args []string
		dir  string
		env  []string
	}{
		{[]string{"init", "--bare"}, tempDir, nil},
		{[]string{"init"}, projPath, nil},
		{[]string{"remote", "add", "origin", remoteURI}, projPath, nil},
		{[]string{"add", "."}, projPath, nil},
		{[]string{"commit", "-m", "test"}, projPath, []string{
			"GIT_COMMITTER_NAME=test",
			"GIT_COMMITTER_EMAIL=test@example.com",
			"GIT_AUTHOR_NAME=test",
			"GIT_AUTHOR_EMAIL=test@example.com",
		}},
	}

	for _, g := range gitCmdList {
		gitCmd := exec.Command(gitExec, g.args...)
		gitCmd.Dir = g.dir

		// if env vars need to be set, we set them alongside existing environment variables
		if g.env != nil {
			gitCmd.Env = append(os.Environ(), g.env...)
		}
		if err := gitCmd.Run(); err != nil {
			t.Fatal(err)
		}
	}

	// return cleanup function
	return func() {
		os.RemoveAll(tempDir)
		// remove .git created by git init
		os.RemoveAll(filepath.Join(projPath, ".git"))
	}
}

func TestRun(t *testing.T) {
	testCases := []struct {
		name   string
		proj   string
		out    string
		expErr error
	}{
		{
			name:   "success",
			proj:   "./testdata/tool/",
			out:    "Go Build: SUCCESS\nGo Test: SUCCESS\nGofmt: SUCCESS\nGit Push: SUCCESS\n",
			expErr: nil,
		},
		{
			name:   "fail",
			proj:   "./testdata/toolErr/",
			out:    "",
			expErr: &stepErr{step: "go build"},
		},
		{
			name:   "failFormat",
			proj:   "./testdata/toolFmtErr/",
			out:    "",
			expErr: &stepErr{step: "go fmt"},
		},
		// {
		// 	name:   "failTimeout",
		// 	proj:   "./testdata/tool",
		// 	out:    "",
		// 	expErr: &stepErr{step: "sleep 11"},
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer

			err := run(tc.proj, &out)

			// when we expecting error to occur
			if tc.expErr != nil {
				if err == nil {
					t.Errorf("Expected error: %q. Got 'nil' instead.", tc.expErr)
				}
				if !errors.Is(err, tc.expErr) {
					t.Errorf("Expected error: %q. Got %q.", tc.expErr, err)
				}
				return //stop execution any further
			}

			// check if unwanted error is there
			if err != nil {
				t.Errorf("Unexpected error: %q", err)
			}

			// validate output
			if out.String() != tc.out {
				t.Errorf("Expected output: %q, Got %q.", tc.out, out.String())
			}
		})
	}
}
