package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
)

func writeStatus(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func serveFile(w http.ResponseWriter, r *http.Request, file string, stat os.FileInfo, extraHeaders map[string]string) {
	if stat == nil {
		s, err := os.Stat(file)
		if err != nil {
			writeStatus(w, http.StatusInternalServerError, fmt.Sprintf("Stat file failed: %v", err))
			return
		}
		stat = s
	}

	io, err := os.Open(file)
	if err != nil {
		writeStatus(w, http.StatusInternalServerError, fmt.Sprintf("Open file failed: %v", err))
		return
	}
	defer io.Close()

	for h, v := range extraHeaders {
		w.Header().Add(h, v)
	}

	http.ServeContent(w, r, path.Base(r.URL.Path), stat.ModTime(), io)
}
