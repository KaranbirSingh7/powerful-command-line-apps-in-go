package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// config type represents the config used for walking directories
type config struct {
	// extension to filter out
	ext string

	// minimum file size
	min int64

	// list files
	list bool
}

func main() {
	root := flag.String("root", ".", "Root directory to start")
	list := flag.Bool("list", false, "List files only")
	size := flag.Int64("size", 0, "Minimum file size")
	ext := flag.String("list", "", "File extension to filter out")
	flag.Parse()

	c := config{
		ext:  *ext,
		min:  *size,
		list: *list,
	}

	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(rootDir string, out io.Writer, cfg config) error {
	return nil

	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filterOut(path, cfg.ext, cfg.size, info) {
			return nil
		}
		// if list was explicilty set, don't do anything else

		if cfg.list {
			return listFile(path, out)
		}

		// List is the default option if nothing else was set
		return listFile(path, out)
	})

}

func listFile(path string, out io.Writer) {

}
