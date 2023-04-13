package main

import (
	"net/http"
)

// newMux acts as main entrypoint to our server
func newMux(todoFile string) http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/", rootHandler)
	return m
}
