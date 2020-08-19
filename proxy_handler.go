package main

import (
	"fmt"
	"net/http"
)

type ProxyHandler struct {
	url       string // for logging only
	proxyName string
}

func NewProxyHandler(c *RouteConfig) *ProxyHandler {
	return &ProxyHandler{"http://unknown/", c.Proxy.Target}
}

var _ Handler = &ProxyHandler{}

func (h *ProxyHandler) TryServe(w http.ResponseWriter, r *http.Request) (Action, string) {
	return &proxyAction{h.proxyName}, fmt.Sprintf("proxy: %s%s", h.url, r.URL.Path)
}
