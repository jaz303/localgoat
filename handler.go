package main

import "net/http"

type Handler interface {
	Start()
	Prefix() string
	TryServe(w http.ResponseWriter, r *http.Request) bool
}
