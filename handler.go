package main

import "net/http"

type Route struct {
	Handler
	Terminal bool
}

type Handler interface {
	Start()
	Prefix() string
	TryServe(w http.ResponseWriter, r *http.Request) (bool, string)
}
