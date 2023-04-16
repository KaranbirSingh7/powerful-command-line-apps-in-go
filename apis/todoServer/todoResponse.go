package main

import (
	"encoding/json"
	"time"

	"github.com/karanbirsingh7/pclaig/todo"
)

// todoResponse
type todoResponse struct {
	Results todo.List `json:"results"`
}

// MarshalJSON - used to marshal struct into json bytes
func (r *todoResponse) MarshalJSON() ([]byte, error) {
	resp := struct {
		Results      todo.List `json:"results"`
		Date         int64     `json:"date"`
		TotalResults int       `json:"total_results"`
	}{
		Results:      r.Results,
		Date:         time.Now().Unix(), //just a timestamp
		TotalResults: len(r.Results),
	}

	return json.Marshal(resp)
}
