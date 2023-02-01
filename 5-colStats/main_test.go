package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkRunAvg(b *testing.B) {
	filenames, err := filepath.Glob("testdata/benchmark/*.csv")
	if err != nil {
		b.Fatal(err)
	}

	// IMPORTANT: ResetTimer ensures that any time used for preparing tests is reset
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := run(filenames, "avg", 2, io.Discard); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkRunMin(b *testing.B) {
	filenames, err := filepath.Glob("testdata/benchmark/*.csv")
	if err != nil {
		b.Fatal(err)
	}

	// IMPORTANT: ResetTimer ensures that any time used for preparing tests is reset
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := run(filenames, "min", 2, io.Discard); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkRunMax(b *testing.B) {
	filenames, err := filepath.Glob("testdata/benchmark/*.csv")
	if err != nil {
		b.Fatal(err)
	}

	// IMPORTANT: ResetTimer ensures that any time used for preparing tests is reset
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := run(filenames, "max", 2, io.Discard); err != nil {
			b.Error(err)
		}
	}
}

func TestRun(t *testing.T) {
	testCases := []struct {
		name   string
		col    int
		op     string
		exp    string
		files  []string
		expErr error
	}{
		{
			name:   "RunAvgFile",
			col:    3,
			op:     "avg",
			exp:    "227.6\n",
			files:  []string{"./testdata/example.csv"},
			expErr: nil,
		},
		{
			name: "RunAvgMultiFiles",
			col:  3,
			op:   "avg",
			exp:  "233.84\n",
			files: []string{
				"./testdata/example.csv",
				"./testdata/example2.csv",
			},
			expErr: nil,
		},
		{
			name: "RunMinmimunFile",
			col:  3,
			op:   "min",
			exp:  "218\n",
			files: []string{
				"./testdata/example.csv",
			},
		},
		{
			name: "RunMinmimunMultiFiles",
			col:  3,
			op:   "min",
			exp:  "218\n",
			files: []string{
				"./testdata/example.csv",
				"./testdata/example2.csv",
			},
		}, {
			name: "RunMaximumFile",
			col:  3,
			op:   "max",
			exp:  "238\n",
			files: []string{
				"./testdata/example.csv",
			},
		},
		{
			name: "RunMaximumMultiFiles",
			col:  3,
			op:   "max",
			exp:  "238\n",
			files: []string{
				"./testdata/example.csv",
				"./testdata/example2.csv",
			},
		},
		{
			name: "RunFailRead",
			col:  2,
			op:   "avg",
			exp:  "",
			files: []string{
				"./testdata/example.csv",
				"./testdata/fake-non-existent-file.csv",
			},
			expErr: os.ErrNotExist, // bcz file#2 doesn't exist
		},
		{
			name: "RunFailColumn",
			col:  0,
			op:   "avg",
			exp:  "",
			files: []string{
				"./testdata/example.csv",
			},
			expErr: ErrInvalidColumn,
		},
		{
			name:   "RunFailNoFiles", // when no file arg is passed by user
			col:    2,
			op:     "avg",
			exp:    "",
			files:  []string{},
			expErr: ErrNoFiles,
		},
		{
			name: "RunFailOperation",
			col:  2,
			op:   "invalid",
			exp:  "",
			files: []string{
				"./testdata/example.csv",
			},
			expErr: ErrInvalidOperation,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var res bytes.Buffer // to capture run function output, could be anything that implements io.Writer
			err := run(tc.files, tc.op, tc.col, &res)

			// check expected errors
			if tc.expErr != nil {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				if !errors.Is(err, tc.expErr) {
					t.Errorf("Expected error %q, got %q instead.", tc.expErr, err)
				}
				return // halt and move onto next testCase
			}

			if err != nil {
				t.Errorf("Unexpected error: %q", err)
			}

			// compare output

			if res.String() != tc.exp {
				t.Errorf("Expected %q, got %q instead", tc.exp, &res)
			}

		})
	}
}
