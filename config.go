package main

type Config struct {
	Address string        `json:"address"`
	Port    int           `json:"port"`
	Routes  []RouteConfig `json:"routes"`
}

type RouteConfig struct {
	Static *StaticRouteConfig `json:"static"`
	Proxy  *ProxyRouteConfig  `json:"proxy"`
}

type StaticRouteConfig struct {
	Path        string            `json:"path"`
	Prefix      string            `json:"prefix"`
	NoCache     bool              `json:"noCache"`
	Exclusive   bool              `json:"exclusive"`
	StripPrefix bool              `json:"stripPrefix"`
	Headers     map[string]string `json:"headers"`
}

type ProxyRouteConfig struct {
	Prefix      string `json:"prefix"`
	Target      string `json:"target"`
	Concurrency int    `json:"concurrency"`
}

func defaultConfig() *Config {
	return &Config{
		Address: "127.0.0.1",
		Port:    8080,
	}
}

func blankConfig() *Config {
	return &Config{
		Address: "",
		Port:    -1,
	}
}
