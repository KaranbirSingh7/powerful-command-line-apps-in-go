package main_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var (
	binName  = "todo"
	fileName = ".todo.json.test"
)

func TestTODOCLI(t *testing.T) {
	task := "test task number 1"
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)

	t.Run("AddNewTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add", task)
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	task2 := "test task number 2"
	t.Run("AddNewTaskFromSTDIN", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add")
		cmdStdin, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}

		// write to stdin pipe and close
		io.WriteString(cmdStdin, task2)
		cmdStdin.Close()

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}

	})

	t.Run("ListTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		expected := fmt.Sprintf("  1: %s\n  2: %s\n", task, task2)
		if string(out) != expected {
			t.Errorf("Got %q, want %q instead\n", string(out), expected)
		}
	})

	t.Run("CompleteTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-complete", "1")
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ListTasksPendingOnly", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list", "-pending")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		// 1 is complete/marked so we shouldn't expect it to be in our output result
		expected := fmt.Sprintf("  2: %s\n", task2)
		if string(out) != expected {
			t.Errorf("Got %q, want %q instead\n", string(out), expected)
		}
	})

	t.Run("ListTasksVerbose", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list", "-verbose")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		expected := fmt.Sprintf("  1: %s\n\tCreated:  2: %s\n", task, task2)
		if !strings.Contains(string(out), "Created") {
			t.Errorf("Got %q, want %q instead\n", string(out), expected)
		}
	})

	t.Run("DeleteTask", func(t *testing.T) {
		itemToDelete := "1"
		cmd := exec.Command(cmdPath, "-delete", itemToDelete)
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})
}

func TestMain(m *testing.M) {
	fmt.Println("Building tool")
	fmt.Println("Setting environment variable TODO_FILENAME=", fileName)

	os.Setenv("TODO_FILENAME", fileName)

	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	// build binary/cli tool
	build := exec.Command("go", "build", "-o", binName)
	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build tool %s: %s", binName, err)
		os.Exit(1)
	}

	fmt.Println("Running tests ...")

	result := m.Run()

	fmt.Println("Cleaning up ... ")

	os.Remove(binName)
	os.Remove(fileName)
	os.Exit(result)
}
