package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

// for mocking exec.CommandContext
func mockCmdContext(ctx context.Context, exe string, args ...string) *exec.Cmd {
	cs := []string{
		"-test.run=TestHelperProcess",
	}
	cs = append(cs, exe)
	cs = append(cs, args...)

	cmd := exec.CommandContext(ctx, os.Args[0], cs...)

	// to ensure that test isn't skipped
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func mockCmdTimeout(ctx context.Context, exe string, args ...string) *exec.Cmd {
	cmd := mockCmdContext(ctx, exe, args...)

	// to ensure that timeout is ON
	cmd.Env = append(cmd.Env, "GO_HELPER_TIMEOUT=1")
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// if timeout is enabled, simulate a timeout of 15s
	if os.Getenv("GO_HELPER_TIMEOUT") == "1" {
		time.Sleep(15 * time.Second)
	}

	if os.Args[2] == "git" {
		fmt.Fprintln(os.Stdout, "Everything up-to-date")
		os.Exit(0)
	}
	os.Exit(1)
}

// setupGit - a test helper func to setup a temporary git repository
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
	// following are NOT required since we now have to ability to mock external commands incld. git
	// _, err := exec.LookPath("git")
	// if err != nil {
	// 	t.Skip("Git not installed, Skipping test.")
	// }
	testCases := []struct {
		name     string
		proj     string
		out      string
		expErr   error
		setupGit bool
		mockCmd  func(ctx context.Context, name string, arg ...string) *exec.Cmd
	}{
		{
			name:     "success",
			proj:     "./testdata/tool/",
			out:      "Go Build: SUCCESS\nGo Test: SUCCESS\nGofmt: SUCCESS\nGit Push: SUCCESS\n",
			expErr:   nil,
			setupGit: true,
			mockCmd:  nil,
		},
		{
			name:     "successMock",
			proj:     "./testdata/tool/",
			out:      "Go Build: SUCCESS\nGo Test: SUCCESS\nGofmt: SUCCESS\nGit Push: SUCCESS\n",
			expErr:   nil,
			setupGit: false,
			mockCmd:  mockCmdContext,
		},
		{
			name:     "fail",
			proj:     "./testdata/toolErr/",
			out:      "",
			expErr:   &stepErr{step: "go build"},
			setupGit: false,
			mockCmd:  nil,
		},
		{
			name:     "failFormat",
			proj:     "./testdata/toolFmtErr/",
			out:      "",
			expErr:   &stepErr{step: "go fmt"},
			setupGit: false,
		},
		{
			name:     "failTimeout",
			proj:     "./testdata/tool",
			out:      "",
			expErr:   context.DeadlineExceeded,
			setupGit: false,
			mockCmd:  mockCmdTimeout,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// setup a temporary git repo if enabled in test case
			if tc.setupGit {
				// check if git cli is present
				_, err := exec.LookPath("git")
				if err != nil {
					t.Skip("Git not installed. Skipping test.")
				}
				cleanup := setupGit(t, tc.proj)
				defer cleanup() // remove temp directories when func is closing
			}

			if tc.mockCmd != nil {
				command = tc.mockCmd
			}

			var out bytes.Buffer

			err := run(tc.proj, &out)

			// when we expecting error to occur
			if tc.expErr != nil {
				if err == nil {
					t.Errorf("Expected error: %q. Got 'nil' instead.", tc.expErr)
					return
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

func TestRunKill(t *testing.T) {
	var testCases = []struct {
		name   string
		proj   string
		sig    syscall.Signal
		expErr error
	}{
		{"SIGINT", "./testdata/tool", syscall.SIGINT, ErrSignal},   // we handle this signal
		{"SIGTERM", "./testdata/tool", syscall.SIGTERM, ErrSignal}, // we handle this signal
		{"SIGQUIT", "./testdata/tool", syscall.SIGQUIT, nil},       // we do NOT handle this anywhere /shrug
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			command = mockCmdTimeout

			errCh := make(chan error)
			ignSigCh := make(chan os.Signal, 1)
			expSigCh := make(chan os.Signal, 1)

			signal.Notify(ignSigCh, syscall.SIGQUIT) // listen for signal
			defer signal.Stop(ignSigCh)              // stop listening for signals

			signal.Notify(expSigCh, tc.sig) // listen
			defer signal.Stop(expSigCh)     // stop listening

			// background test our func
			go func() {
				errCh <- run(tc.proj, io.Discard)
			}()

			// use kill signal after 2 seconds has passed
			go func() {
				time.Sleep(2 * time.Second)
				syscall.Kill(syscall.Getpid(), tc.sig)
			}()

			select {
			// if we get an error on channel
			case err := <-errCh:
				if err == nil {
					t.Errorf("Expected error. Got nil instead.")
					return
				}
				if !errors.Is(err, tc.expErr) {
					t.Errorf("Expected error: %q. Got %q", tc.expErr, err)
				}
				// check what kind of error we expecting
				select {
				case rec := <-expSigCh:
					if rec != tc.sig {
						t.Errorf("Expected signal %q, got %q", tc.sig, rec)
					}
				default:
					t.Errorf("Signal not received")
				}
			case <-ignSigCh:
			}
		})
	}
}
