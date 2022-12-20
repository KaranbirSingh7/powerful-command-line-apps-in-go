package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// item struct represents a toDo item
type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

// List represent list of all toDo items
type List []item

// String prints out formatted list
func (l *List) String() string {
	formatted := ""
	verboseOutput := false

	if os.Getenv("DEBUG") != "" {
		verboseOutput = true
	}

	for k, t := range *l {
		prefix := "  "
		if os.Getenv("TODO_SHOW") == "pending" && t.Done {
			continue
		}
		if t.Done {
			prefix = "X "
		}
		formatted += fmt.Sprintf("%s%d: %s\n", prefix, k+1, t.Task)

		if verboseOutput {
			formatted += fmt.Sprintf("\tCreated: %s\n", t.CreatedAt)
		}
	}
	return formatted
}

// Add takes care of adding new todo item to list
func (l *List) Add(task string) {
	t := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{}, //empty time 0000-0000:0000
	}
	// append new item to existing list (modifying underlying pointer value)
	*l = append(*l, t)
}

// Complete methods marks a todo item as complete by setting index=True
func (l *List) Complete(i int) error {
	ls := *l

	// sanity check the value provided
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("item %d does not exist", i)
	}

	ls[i-1].Done = true
	ls[i-1].CompletedAt = time.Now()

	return nil
}

// Delete deletes an item from list
func (l *List) Delete(i int) error {
	ls := *l

	// sanity check the value provided
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("item %d does not exist", i)
	}

	*l = append(ls[:i-1], ls[i:]...)

	return nil
}

// Save writes list to a JSON file
func (l *List) Save(filename string) error {
	json, err := json.Marshal(l) //convert list to json bytes
	if err != nil {
		return err
	}
	return os.WriteFile(filename, json, 0644)
}

// Get reads JSON from a file into list
func (l *List) Get(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		// if file not found, there is no error to be returned
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	// file exists but empty
	if len(file) == 0 {
		return nil
	}

	return json.Unmarshal(file, l)

}
