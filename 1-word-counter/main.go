package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	CLINAME = "word counter"
	VERSION = "0.0.1"
)

func main() {
	linesFlag := flag.Bool("l", false, "count lines")
	bytesFlag := flag.Bool("b", false, "count bytes")
	inputFile := flag.String("f", "", "file to read input from")
	flag.Parse()

	// log.Printf("Starting program - %s [v%s]\n", CLINAME, VERSION)
	// log.Printf("Word count: %d\n", count(os.Stdin))

	if *inputFile != "" {
		f, err := os.Open(*inputFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		fmt.Println(count(f, *linesFlag, *bytesFlag))
		return //halt
	}

	fmt.Println(count(os.Stdin, *linesFlag, *bytesFlag))

}

// count takes string as input and return the count of works in a string
func count(r io.Reader, countLines bool, countBytes bool) int {
	// a scanner is used to read text from a Reader (such as files or STDIN)
	scanner := bufio.NewScanner(r) // init a new reader on our reader interface (could be a file or in this case STDIN)

	// if countLine flag is not set, we wanna split words, if passed/set we wanna use scanLines(default)
	// define the scanner split type (default is split by lines)
	// following result "a simple string" -> ['a', 'simple', 'string']
	if !countLines {
		scanner.Split(bufio.ScanWords)
	}

	// define a counter variable
	var wordCount int
	var bytesCount int

	if countBytes {
		// find the # of bytes
		for scanner.Scan() {
			bytesCount += len(scanner.Bytes())
		}
		return bytesCount
	} else {
		for scanner.Scan() {
			// for every word count, increment the counter
			wordCount++
		}
	}
	return wordCount
}
