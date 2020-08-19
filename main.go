package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"

	"github.com/felixge/httpsnoop"
)

type proxyTable map[string]*httputil.ReverseProxy

type route struct {
	Handler
	Prefix   string
	Terminal bool
}

func getRoutes(cfg *Config, proxies proxyTable) []route {
	out := make([]route, 0)
	for _, r := range cfg.Routes {
		var h Handler
		if r.Static != nil {
			h = NewStaticHandler(&r)
		} else if r.Proxy != nil {
			h = NewProxyHandler(&r)
		}
		out = append(out, route{h, r.Prefix, r.Terminal})
	}
	return out
}

var (
	messages    map[*http.Request]string = make(map[*http.Request]string)
	messageLock sync.Mutex
)

func buildProxyTable(cfg *Config) proxyTable {
	proxies := make(proxyTable)
	for name, pc := range cfg.Proxies {
		url, err := url.Parse(pc.Host)
		if err != nil {
			// TODO: handle error
		}
		proxies[name] = httputil.NewSingleHostReverseProxy(url)
	}
	return proxies
}

func setMessage(r *http.Request, msg string) {
	messageLock.Lock()
	defer messageLock.Unlock()
	messages[r] = msg
}

func getMessage(r *http.Request) string {
	messageLock.Lock()
	defer messageLock.Unlock()
	msg, ok := messages[r]
	if ok {
		delete(messages, r)
	}
	return msg
}

func main() {
	cfg := getConfiguration()
	proxies := buildProxyTable(cfg)
	routes := getRoutes(cfg, proxies)
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var action Action
		var message string

		for _, route := range routes {
			if !strings.HasPrefix(r.URL.Path, route.Prefix) {
				continue
			}
			action, message = route.TryServe(w, r)
			if action != nil || route.Terminal {
				break
			}
		}

		if action == nil {
			action = &messageAction{http.StatusNotFound, "Not Found"}
			message = "no matching route"
		}

		switch m := action.(type) {
		case *messageAction:
			w.WriteHeader(m.StatusCode)
			w.Write([]byte(m.Message))
			break
		case *serveFileAction:
			serveFile(w, r, m.TargetFile, m.Stat, m.ExtraHeaders)
			break
		case *proxyAction:
			rp, ok := proxies[m.ProxyName]
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Unknown proxy name: %s", m.ProxyName)))
			} else {
				rp.ServeHTTP(w, r)
			}
			break
		}

		setMessage(r, message)
	})

	addr := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
	fmt.Printf("localgoat starting, listening on %s\n", addr)

	http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := httpsnoop.CaptureMetrics(mux, w, r)
		fmt.Printf("%d %s %s (%s)\n", m.Code, r.Method, r.URL.Path, getMessage(r))
	}))
}
