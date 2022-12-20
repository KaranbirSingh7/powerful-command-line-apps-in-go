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
	size int64

	// list files
	list bool

	// delete files if enabled
	del bool
}

func main() {
	root := flag.String("root", ".", "root directory to start")
	list := flag.Bool("list", false, "list files only")
	size := flag.Int64("size", 0, "minimum file size")
	ext := flag.String("ext", "", "file extension to filter out")
	del := flag.Bool("del", false, "delete matching files")
	flag.Parse()

	c := config{
		ext:  *ext,
		size: *size,
		list: *list,
		del:  *del,
	}

	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(rootDir string, out io.Writer, cfg config) error {
	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
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

		// if delete flag is passed
		if cfg.del {
			return delFile(path)
		}

		// List is the default option if nothing else was set
		return listFile(path, out)
	})

}
