package main

import "net/http"

type Handler interface {
	Start()
	TryServe(w http.ResponseWriter, r *http.Request) bool
}
