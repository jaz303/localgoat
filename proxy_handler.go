package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type ProxyHandler struct {
	config *ProxyRouteConfig
	proxy  *httputil.ReverseProxy
}

var _ Handler = &ProxyHandler{}

func NewProxyHandler(c *ProxyRouteConfig) *ProxyHandler {
	if c.Concurrency < 1 {
		c.Concurrency = 1
	}
	return &ProxyHandler{
		config: c,
	}
}

func (h *ProxyHandler) Start() {
	target, _ := url.Parse(h.config.Target)
	h.proxy = httputil.NewSingleHostReverseProxy(target)
}

func (h *ProxyHandler) Prefix() string {
	return h.config.Prefix
}

func (h *ProxyHandler) TryServe(w http.ResponseWriter, r *http.Request) (bool, string) {
	if !strings.HasPrefix(r.URL.Path, h.config.Prefix) {
		return false, ""
	}

	h.proxy.ServeHTTP(w, r)
	return true, fmt.Sprintf("proxy to %s%s", h.config.Target, r.URL.Path)
}
