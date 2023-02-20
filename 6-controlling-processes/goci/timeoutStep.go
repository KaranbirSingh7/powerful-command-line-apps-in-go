package main

import (
	"context"
	"os/exec"
	"time"
)

type timeoutStep struct {
	step
	timeout time.Duration
}

// newTimeoutStep - constructor for creating a new timeout object
func newTimeoutStep(name, exe, message, proj string, args []string, timeout time.Duration) timeoutStep {
	s := timeoutStep{}
	s.step = newStep(name, exe, message, proj, args)
	s.timeout = timeout

	// if nothing provided, then set 30seconds as default timeout.
	if s.timeout == 0 {
		s.timeout = 30 * time.Second
	}
	return s
}

// assign func type to command variable - this will populate command var with func defaults
// command is a package variable - unexported
var command = exec.CommandContext

func (s timeoutStep) execute() (string, error) {
	// init our context
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	cmd := command(ctx, s.exe, s.args...)
	cmd.Dir = s.proj

	if err := cmd.Run(); err != nil {
		// if error points to timeout
		if ctx.Err() == context.DeadlineExceeded {
			return "", &stepErr{
				step:  s.name,
				msg:   "failed time out",
				cause: context.DeadlineExceeded,
			}
		}
		return "", &stepErr{
			step:  s.name,
			msg:   "failed to execute",
			cause: err,
		}
	}

	// message provided by user  + no error :-)
	return s.message, nil
}
