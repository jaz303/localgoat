package main

type Config struct {
	Address string        `json:"address"`
	Port    int           `json:"port"`
	Routes  []RouteConfig `json:"routes"`
}

type RouteConfig struct {
	Static   *StaticRouteConfig `json:"static"`
	Proxy    *ProxyRouteConfig  `json:"proxy"`
	Terminal bool               `json:"terminal"`
}

type StaticRouteConfig struct {
	Path        string            `json:"path"`
	Prefix      string            `json:"prefix"`
	NoCache     bool              `json:"noCache"`
	StripPrefix bool              `json:"stripPrefix"`
	Headers     map[string]string `json:"headers"`
}

type ProxyRouteConfig struct {
	Prefix      string `json:"prefix"`
	Target      string `json:"target"`
	Concurrency int    `json:"concurrency"`
}

func blankConfig() *Config {
	return &Config{
		Address: "",
		Port:    -1,
	}
}
