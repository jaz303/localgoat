package main

import "net/http"

func writeHeaders(w http.ResponseWriter, hs map[string]string) {
	for h, v := range hs {
		w.Header().Add(h, v)
	}
}

func thwartCache(w http.ResponseWriter) {
	writeHeaders(w, map[string]string{
		"Cache-Control": "no-cache, no-store, must-revalidate",
		"Pragma":        "no-cache",
		"Expires":       "0",
	})
}
