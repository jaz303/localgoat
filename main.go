package main

import (
	"fmt"
	"net/http"
	"strings"
)

func getRoutes(cfg *Config) []Handler {
	out := make([]Handler, 0)
	for _, r := range cfg.Routes {
		if r.Static != nil {
			out = append(out, NewStaticHandler(r.Static))
		} else if r.Proxy != nil {
			out = append(out, NewProxyHandler(r.Proxy))
		}
	}
	return out
}

func main() {
	cfg := getConfiguration()

	routes := getRoutes(cfg)
	for _, route := range routes {
		route.Start()
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for _, route := range routes {
			if strings.HasPrefix(r.URL.Path, route.Prefix()) && route.TryServe(w, r) {
				return
			}
		}
		writeNotFound(w)
	})

	http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Address, cfg.Port), mux)
}
