package main

import (
	"errors"
	"fmt"
)

var (
	ErrValidation = errors.New("validation failed")
	ErrSignal     = errors.New("received signal")
)

// stepErr - used to carry information about an error
type stepErr struct {
	step  string
	msg   string
	cause error
}

// Error - used to satisfy errors interface
func (s *stepErr) Error() string {
	return fmt.Sprintf("Step: %q: %s: Cause: %v", s.step, s.msg, s.cause)
}

// Is - check if error is of type stepErr
func (s *stepErr) Is(target error) bool {
	t, ok := target.(*stepErr)
	if !ok {
		return false
	}
	return t.step == s.step
}

// Unwrap - unwraps the error
func (s *stepErr) Unwrap() error {
	return s.cause
}
