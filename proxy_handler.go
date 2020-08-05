package main

import "net/http"

type ProxyHandler struct {
	config *ProxyRouteConfig
}

var _ Handler = &ProxyHandler{}

func NewProxyHandler(c *ProxyRouteConfig) *ProxyHandler {
	return &ProxyHandler{c}
}

func (h *ProxyHandler) TryServe(w http.ResponseWriter, r *http.Request) bool {
	return false
}
