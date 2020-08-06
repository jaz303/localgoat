package main

import (
	"net/http"
	"os"
	"path"
	"strings"
)

type StaticHandler struct {
	config *StaticRouteConfig
}

var _ Handler = &StaticHandler{}

func NewStaticHandler(c *StaticRouteConfig) *StaticHandler {
	return &StaticHandler{c}
}

func (h *StaticHandler) TryServe(w http.ResponseWriter, r *http.Request) bool {
	if !strings.HasPrefix(r.URL.Path, h.config.Prefix) {
		return false
	}

	targetFile := path.Join(h.config.Path, r.URL.Path)

	stat, err := os.Stat(targetFile)
	if err != nil {
		return false
	}

	io, err := os.Open(targetFile)
	if err != nil {
		return false
	}
	defer io.Close()

	writeHeaders(w, h.config.Headers)

	if h.config.NoCache {
		thwartCache(w)
	}

	http.ServeContent(w, r, path.Base(r.URL.Path), stat.ModTime(), io)

	return true
}
