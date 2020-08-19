package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/alecthomas/kong"
)

type CLIOptions struct {
	Version       bool     `name:"version" help:"Print version and exit."`
	ConfigFile    string   `name:"config" short:"c" placeholder:"CONFIG-FILE" help:"Config file. If specified, all other configuration options will be ignored."`
	Address       string   `name:"address" short:"a" default:"127.0.0.1" help:"Listen address."`
	Port          int      `name:"port" short:"p" default:"8080" help:"Listen port."`
	DisableCache  bool     `name:"disable-cache" short:"C" default:"0" help:"Inject headers to thwart HTTP caching (static content only)."`
	StaticPrefix  string   `name:"static-prefix" default:"/" help:"URL prefix for static files."`
	NoStripPrefix bool     `name:"no-strip-prefix" short:"S" default:"0" help:"Don't strip prefix when resolving static files."`
	Args          []string `arg optional help:"When no config file is specified trailing arguments invoke common use-cases; behaviour is determined by argument count. (0): serve static content from current directory. (1): serve static content from the specified directory. (2): proxy requests with prefix arg0 to target arg1. (3): serve static content from arg0, and proxy requests with prefix arg1 to target arg2."`
}

func loadConfigFromFile(file string) *Config {
	f, _ := os.Open(file)
	data, _ := ioutil.ReadAll(f)
	cfg := Config{}
	json.Unmarshal(data, &cfg)
	return &cfg
}

func addDefaultStaticRoute(cfg *Config, dir string, opts *CLIOptions) {
	cfg.Routes = append(cfg.Routes, RouteConfig{
		Static: &StaticRouteConfig{
			Path:        dir,
			NoCache:     opts.DisableCache,
			StripPrefix: !opts.NoStripPrefix,
			Headers:     make(map[string]string),
		},
		Prefix: opts.StaticPrefix,
	})
}

func addDefaultProxyRoute(cfg *Config, prefix string, target string, opts *CLIOptions) {
	cfg.Proxies["default"] = ProxyConfig{
		Host: target,
	}
	cfg.Routes = append(cfg.Routes, RouteConfig{
		Proxy: &ProxyRouteConfig{
			Target: "default",
		},
		Prefix: prefix,
	})
}

func getConfiguration() *Config {
	var CLI CLIOptions

	kong.Parse(&CLI)

	if CLI.Version {
		fmt.Printf("localgoat v%s\n", localgoatVersion)
		os.Exit(0)
	}

	if len(CLI.ConfigFile) > 0 {
		return loadConfigFromFile(CLI.ConfigFile)
	}

	cfg := Config{}
	cfg.Address = CLI.Address
	cfg.Port = CLI.Port
	cfg.Proxies = make(map[string]ProxyConfig)

	switch len(CLI.Args) {
	case 0:
		addDefaultStaticRoute(&cfg, ".", &CLI)
	case 1:
		addDefaultStaticRoute(&cfg, CLI.Args[0], &CLI)
	case 2:
		addDefaultProxyRoute(&cfg, CLI.Args[0], CLI.Args[1], &CLI)
	case 3:
		addDefaultStaticRoute(&cfg, CLI.Args[0], &CLI)
		addDefaultProxyRoute(&cfg, CLI.Args[1], CLI.Args[2], &CLI)
	default:
	}

	return &cfg
}
