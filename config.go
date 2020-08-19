package main

type Config struct {
	Address string                 `json:"address"`
	Port    int                    `json:"port"`
	Proxies map[string]ProxyConfig `json:"proxies"`
	Routes  []RouteConfig          `json:"routes"`
}

type ProxyConfig struct {
	Host string `json:"host"`
}

type RouteConfig struct {
	Static   *StaticRouteConfig `json:"static"`
	Proxy    *ProxyRouteConfig  `json:"proxy"`
	Prefix   string             `json:"prefix"`
	Terminal bool               `json:"terminal"`
}

type StaticRouteConfig struct {
	Path        string            `json:"path"`
	NoCache     bool              `json:"noCache"`
	StripPrefix bool              `json:"stripPrefix"`
	Headers     map[string]string `json:"headers"`
}

type ProxyRouteConfig struct {
	Target string `json:"target"`
}

func blankConfig() *Config {
	return &Config{
		Address: "",
		Port:    -1,
	}
}
