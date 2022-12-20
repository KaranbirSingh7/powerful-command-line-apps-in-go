package main

import (
	"bytes"
	"os"
	"testing"
)

const (
	testFile = "testdata/log.txt"
)

// TestCount test the function that count words
func TestCountWords(t *testing.T) {
	words := bytes.NewBufferString("a simple string\n")
	want := 3
	got := count(words, false, false)
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func TestCountLines(t *testing.T) {
	b := bytes.NewBufferString("word1 word2\n word3\n word4")
	got := count(b, true, false)
	want := 3
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func TestCountBytes(t *testing.T) {
	b := bytes.NewBufferString("a \n simple \n byte")
	got := count(b, false, true)
	want := 11
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func TestCountWordsFromFile(t *testing.T) {
	// open test file for reading
	f, err := os.Open(testFile)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	want := 9
	got := count(f, false, false)

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func TestCountLinesFromFile(t *testing.T) {
	// open test file for reading
	f, err := os.Open(testFile)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	want := 3
	got := count(f, true, false)

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func TestCountTable(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{input: "a simple line", want: 3},
		{input: "a line with five words", want: 5},
		{input: "a line with six s words", want: 6},
		{input: "a line with seven boo n words", want: 7},
	}

	for _, tc := range tests {
		got := count(bytes.NewBufferString(tc.input), false, false)
		want := tc.want
		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	}
}
