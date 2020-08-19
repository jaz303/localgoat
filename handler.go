package main

import "net/http"

type Handler interface {
	TryServe(w http.ResponseWriter, r *http.Request) (Action, string)
}
