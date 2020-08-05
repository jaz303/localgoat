package main

import "net/http"

type StaticHandler struct {
	config *StaticRouteConfig
}

var _ Handler = &StaticHandler{}

func NewStaticHandler(c *StaticRouteConfig) *StaticHandler {
	return &StaticHandler{c}
}

func (h *StaticHandler) TryServe(w http.ResponseWriter, r *http.Request) bool {
	return false
}
