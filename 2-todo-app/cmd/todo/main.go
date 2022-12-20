package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/karanbirsingh7/pclaig/todo"
)

var todoFileName = ".todo.json"

func main() {

	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	// cli flags
	add := flag.Bool("add", false, "Add task to ToDo list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be mark as completed")
	delete := flag.Int("delete", 0, "Item to delete from list")
	verbose := flag.Bool("verbose", false, "Verbose output when listing tasks")
	pending := flag.Bool("pending", false, "Show only pending items")

	flag.Usage = func() {
		fmt.Println("My TODO CLI")
		flag.PrintDefaults()
	}

	flag.Parse()

	l := &todo.List{}

	// if any issues with reading file, print to STDERR and exit
	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if *verbose {
		os.Setenv("DEBUG", "true")
	}

	if *pending {
		os.Setenv("TODO_SHOW", "pending") // show only pending items
	}

	switch {
	case *list:
		// list flag means list all items
		fmt.Print(l)
	case *delete > 0:
		// delete the item
		if err := l.Delete(*delete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// save new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *complete > 0:
		// complete given item
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// save new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case *add:
		// when any arguments are provided, they will be used as new task
		t, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// add new task
		l.Add(t)

		// save list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	// default print flags
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}

// getTask decides where to get new task from, could be STDIN or arguments
func getTask(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	s := bufio.NewScanner(r)
	s.Scan()

	if err := s.Err(); err != nil {
		return "", err
	}

	if len(s.Text()) == 0 {
		return "", fmt.Errorf("task cannot be blank")
	}

	return s.Text(), nil

}
