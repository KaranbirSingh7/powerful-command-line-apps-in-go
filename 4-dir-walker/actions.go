package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// filterOut is responsbile for checking if path ...
func filterOut(path string, ext string, minSize int64, info os.FileInfo) bool {
	if info.IsDir() || info.Size() < minSize {
		return true
	}

	if ext != "" && filepath.Ext(path) != ext {
		return true
	}
	return false
}

// listFile is used to print output of path
func listFile(path string, out io.Writer) error {
	_, err := fmt.Fprintln(out, path)
	return err
}

// delFile takes care of deleting a file
func delFile(path string) error {
	return os.Remove(path)
}
