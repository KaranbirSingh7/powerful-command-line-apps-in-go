package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var opFunc statsFunc

func main() {

	// verify and parse arguments
	op := flag.String("op", "sum", "operation to run on selected column")
	column := flag.Int("col", 1, "CSV column to run operation on")

	flag.Parse()

	if err := run(flag.Args(), *op, *column, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filenames []string, op string, column int, out io.Writer) error {
	// no file meaning nothing to run on
	if len(filenames) == 0 {
		return ErrNoFiles
	}

	// if user provided a negative number
	if column < 1 {
		return fmt.Errorf("%w: %d", ErrInvalidColumn, column)
	}

	switch op {
	case "sum":
		opFunc = sum
	case "avg":
		opFunc = avg
	default:
		return fmt.Errorf("%w: %s", ErrInvalidOperation, op)
	}

	consolidate := make([]float64, 0)

	// loop through all filenames
	for _, fname := range filenames {
		// Open the file for reading
		f, err := os.Open(fname)
		if err != nil {
			return fmt.Errorf("cannot open file: %w", err)
		}

		// Parse the CSV into a slice of float64 number
		data, err := csv2float(f, column)
		if err != nil {
			return err
		}

		// close file to release resources
		if err := f.Close(); err != nil {
			return err
		}
		// append slices to master
		consolidate = append(consolidate, data...)
	}
	_, err := fmt.Fprintln(out, opFunc(consolidate))
	return err
}
