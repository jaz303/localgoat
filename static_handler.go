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
	prefix string
	config *StaticRouteConfig
	spa    bool
}

var _ Handler = &StaticHandler{}

var errIllegalStaticPath = errors.New("illegal path")

func NewStaticHandler(c *RouteConfig) *StaticHandler {
	stat, err := os.Stat(c.Static.Path)
	if err != nil {
		// TODO: error
	}
	return &StaticHandler{c.Prefix, c.Static, !stat.IsDir()}
}

func (h *StaticHandler) TryServe(w http.ResponseWriter, r *http.Request) (Action, string) {
	var targetFile string
	if h.spa {
		targetFile = h.config.Path
	} else {
		resolved, err := h.resolvePath(r.URL.Path)
		if err == errIllegalStaticPath {
			return NotFound, "static: illegal path"
		} else if err != nil {
			return InternalServerError, fmt.Sprintf("static: %v", err)
		}
		targetFile = resolved
	}

	stat, err := os.Stat(targetFile)
	if err != nil {
		return nil, ""
	}

	// TODO: implement directory indexes and index files
	if stat.IsDir() {
		return nil, ""
	}

	action := &serveFileAction{
		targetFile,
		stat,
		make(map[string]string),
	}

	if h.config.NoCache {
		action.ExtraHeaders["Cache-Control"] = "no-cache, no-store, must-revalidate"
		action.ExtraHeaders["Pragma"] = "no-cache"
		action.ExtraHeaders["Expires"] = "0"
	}

	for h, v := range h.config.Headers {
		action.ExtraHeaders[h] = v
	}

	return action, fmt.Sprintf("static: serve %s", targetFile)
}

func (h *StaticHandler) resolvePath(p string) (string, error) {
	if strings.Contains(p, "..") {
		return "", errIllegalStaticPath
	}

	if h.config.StripPrefix {
		p = p[len(h.prefix):]
	}

	return path.Join(h.config.Path, p), nil
}
