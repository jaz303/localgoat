package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
)

type StaticHandler struct {
	config *StaticRouteConfig
}

var _ Handler = &StaticHandler{}

var errIllegalStaticPath = errors.New("illegal path")

func NewStaticHandler(c *StaticRouteConfig) *StaticHandler {
	return &StaticHandler{c}
}

func (h *StaticHandler) Start() {

}

func (h *StaticHandler) Prefix() string {
	return h.config.Prefix
}

func (h *StaticHandler) TryServe(w http.ResponseWriter, r *http.Request) (bool, string) {
	targetFile, err := h.resolvePath(r.URL.Path)
	if err == errIllegalStaticPath {
		writeNotFound(w)
		return true, "static: illegal path"
	} else if err != nil {
		writeInteralServerError(w)
		return true, fmt.Sprintf("static: %v", err)
	}

	stat, err := os.Stat(targetFile)
	if err != nil {
		if h.config.Exclusive {
			writeNotFound(w)
			return true, "static: file not found"
		} else {
			return false, ""
		}
	}

	// TODO: implement directory indexes and index files
	if stat.IsDir() {
		writeNotFound(w)
		return true, "static: can't serve directory"
	}

	io, err := os.Open(targetFile)
	if err != nil {
		return false, ""
	}
	defer io.Close()

	writeHeaders(w, h.config.Headers)

	if h.config.NoCache {
		thwartCache(w)
	}

	http.ServeContent(w, r, path.Base(r.URL.Path), stat.ModTime(), io)

	return true, fmt.Sprintf("static: serve %s", targetFile)
}

func (h *StaticHandler) resolvePath(p string) (string, error) {
	if strings.Contains(p, "..") {
		return "", errIllegalStaticPath
	}

	if h.config.StripPrefix {
		p = p[len(h.Prefix()):]
	}

	return path.Join(h.config.Path, p), nil
}
