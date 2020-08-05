package main

import (
	"fmt"
	"net/http"
)

var routes []Handler

func main() {
	cfg := defaultConfig()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for _, route := range routes {
			if route.TryServe(w, r) {
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	})

	http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Address, cfg.Port), mux)
}
