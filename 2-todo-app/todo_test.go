package todo_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/karanbirsingh7/pclaig/todo"
)

// TestAdd tests the Add method of the list type
func TestAdd(t *testing.T) {
	l := todo.List{}
	taskName := "NewTask"
	l.Add(taskName)

	if l[0].Task != taskName {
		t.Errorf("Got %q, want %q", l[0].Task, taskName)
	}
}

// TestComplete tests the complete method of the List type
func TestComplete(t *testing.T) {
	l := todo.List{}
	taskName := "New Task"
	l.Add(taskName)

	if l[0].Task != taskName {
		t.Errorf("Got %q, want %q", l[0].Task, taskName)
	}

	if l[0].Done {
		t.Errorf("New task should not be completed")
	}

	// mark task as finished
	l.Complete(1)

	if !l[0].Done {
		t.Errorf("New task should be completed")
	}

	if l[0].Done != true {
		t.Errorf("Got %q, want %q", l[0].Task, taskName)
	}
}

// TestDelete tests the delete functionality of list method
func TestDelete(t *testing.T) {
	l := todo.List{}

	// create 3 tasks
	tasks := []string{
		"New Task 1",
		"New Task 2",
		"New Task 3",
	}

	// fill up our list
	for _, v := range tasks {
		l.Add(v)
	}

	if l[0].Task != tasks[0] {
		t.Errorf("Got %q, want %q", l[0].Task, tasks[0])
	}

	// delete item number 2
	l.Delete(2)

	if len(l) != 2 {
		t.Errorf("Expected list length: %d, got %d instead", 2, len(l))
	}

	if l[1].Task != tasks[2] {
		t.Errorf("Got %q, want %q", l[1].Task, tasks[2])
	}

}

// TestSaveGet writes to a tempFile and then reads it back
func TestSaveGet(t *testing.T) {
	l1 := todo.List{}
	l2 := todo.List{}

	taskName := "New Task"
	l1.Add(taskName)

	if l1[0].Task != taskName {
		t.Errorf("Got %q, want %q", l1[0].Task, taskName)
	}

	tf, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("Error creating temporary file: %s", err)
	}
	// remove file on exit
	defer os.Remove(tf.Name())

	// SAVE to a file
	if err := l1.Save(tf.Name()); err != nil {
		t.Fatalf("Error saving list to file: %s", err)
	}

	// Get/READs from a file
	if err := l2.Get(tf.Name()); err != nil {
		t.Fatalf("Error reading from file: %s", err)
	}

	if l1[0].Task != l2[0].Task {
		t.Errorf("Task %q should be match %q task", l1[0].Task, l2[0].Task)
	}

}
