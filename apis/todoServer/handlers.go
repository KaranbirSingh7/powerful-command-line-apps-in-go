package main

import (
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r) // automatic 404
		return
	}
	content := "There's an API here"
	replyTextContext(w, r, http.StatusOK, content)
}

func replyTextContext(w http.ResponseWriter, r *http.Request, status int, content string) {
	w.Header().Set("Content-Type", "text/plain") // set header
	w.WriteHeader(status)                        // write header
	w.Write([]byte(content))                     // write content to response stream
}
