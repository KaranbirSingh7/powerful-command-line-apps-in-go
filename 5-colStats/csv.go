package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

// statsFunc defines a generic statistical function
type statsFunc func(data []float64) float64

// sum - calculate total/sum from list of floats
func sum(data []float64) float64 {
	var result float64
	for _, v := range data {
		result += v
	}
	return result
}

// avg - calculate average from list of floats
func avg(data []float64) float64 {
	return sum(data) / float64(len(data))
}

// csv2float - reads csv data from reader and return slice of floats
func csv2float(r io.Reader, column int) ([]float64, error) {
	// new CSV reader used to read data from CSV files
	cr := csv.NewReader(r)

	// from user perpective, column counts starts with 1 (not 0), jokes on you developers
	// user inputs "1" -> we read it as "0"
	column--

	// read all data in CSV
	allData, err := cr.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Cannot read data from file: %w", err)
	}

	var data []float64

	for i, row := range allData {
		// if header aka position 0, we skip because first row is column names
		if i == 0 {
			continue
		}
		if len(row) <= column {
			return nil, fmt.Errorf("%w: File has only %d columns", ErrInvalidColumn, len(row))
		}

		// try to convert data read into a float number
		v, err := strconv.ParseFloat(row[column], 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrNotNumber, err)
		}

		// if all good, append to our resultingList
		data = append(data, v)
	}
	return data, nil
}
