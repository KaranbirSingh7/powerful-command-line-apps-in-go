package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
)

var opFunc statsFunc

func main() {
	// verify flags and parse arguments
	op := flag.String("op", "sum", "operation to run on selected column, valid options are 'sum', 'avg', 'min' and 'max'")
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
	case "min":
		opFunc = min
	case "max":
		opFunc = max
	default:
		return fmt.Errorf("%w: %s", ErrInvalidOperation, op)
	}

	consolidate := make([]float64, 0)

	// create new channels to receive result, erros and done
	// all channels are un-buffered meaning one in and one out
	resCh := make(chan []float64)
	errCh := make(chan error)
	doneCh := make(chan struct{})
	filesCh := make(chan string)

	// loop through all files sending them through channel
	// so that each one is processed when worker is available
	go func() {
		defer close(filesCh)
		for _, fname := range filenames {
			filesCh <- fname
		}
	}()

	wg := sync.WaitGroup{}

	// following for condition ensures that we are only running N number of goroutines at a time.
	// where N == number of CPUs
	for i := 0; i < runtime.NumCPU(); i++ {
		// increase wg counter
		wg.Add(1)

		go func() {
			defer wg.Done() // decrease counter

			for fname := range filesCh { // for listens to channel until its closed
				// Open the file for reading
				f, err := os.Open(fname)
				if err != nil {
					errCh <- fmt.Errorf("cannot open file: %w", err)
					return // halt
				}

				// Parse the CSV into a slice of float64 number
				data, err := csv2float(f, column)
				if err != nil {
					errCh <- err
				}

				// close file to release resources
				if err := f.Close(); err != nil {
					errCh <- err
				}

				// append slices to master
				// appending to array could cause a race condition to appear since multiple goroutines can be accessing/modifying same variable
				// consolidate = append(consolidate, data...)

				resCh <- data // let responseChannel handle stream of data

			}
		}()
	}
	go func() {
		wg.Wait()     // wait for goroutines to finish
		close(doneCh) // close the doneCh signifying that work has been completed.
	}()

	for {
		select { //blocking statement that watches multiple channels
		case err := <-errCh:
			return err
		case data := <-resCh:
			consolidate = append(consolidate, data...)
		case <-doneCh:
			_, err := fmt.Fprintln(out, opFunc(consolidate))
			return err
		}
	}
	// no return here - for:select:case takes care of it
}
