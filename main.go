package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func loadConfig() *Config {
	f, _ := os.Open("config.example.json")
	data, _ := ioutil.ReadAll(f)
	cfg := Config{}
	json.Unmarshal(data, &cfg)
	return &cfg
}

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
	cfg := loadConfig()

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
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	})

	http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Address, cfg.Port), mux)
}
