package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
)

func loadConfigFromFile(file string) *Config {
	f, _ := os.Open(file)
	data, _ := ioutil.ReadAll(f)
	cfg := Config{}
	json.Unmarshal(data, &cfg)
	return &cfg
}

func usage(printHelp bool, exitCode int) {
	println("Usage: lg [options] [<proxy-prefix> <proxy-target>]")
	if printHelp {
		println("")
		flag.PrintDefaults()
	}
	os.Exit(exitCode)
}

func getConfiguration() *Config {
	var configFile string

	cliConfig := &Config{
		Routes: []RouteConfig{
			{Static: &StaticRouteConfig{}},
		},
	}

	help := flag.Bool("h", false, "Display help")

	flag.StringVar(&configFile, "c", "", "Config file (if set, all other command line options will be ignored)")
	flag.StringVar(&cliConfig.Address, "a", "127.0.0.1", "Listen address")
	flag.IntVar(&cliConfig.Port, "p", 8080, "Listen port")
	flag.StringVar(&cliConfig.Routes[0].Static.Path, "d", ".", "Static directory")

	flag.StringVar(&cliConfig.Routes[0].Static.Prefix, "static-prefix", "/", "Static prefix")
	flag.BoolVar(&cliConfig.Routes[0].Static.NoCache, "disable-cache", false, "Disable caching for static files")
	flag.BoolVar(&cliConfig.Routes[0].Static.StripPrefix, "strip-prefix", true, "Strip prefix for static files")

	flag.Parse()

	if *help {
		usage(true, 0)
	}

	if len(configFile) > 0 {
		return loadConfigFromFile(configFile)
	}

	if flag.NArg() == 2 {
		cliConfig.Routes = append(cliConfig.Routes, RouteConfig{Proxy: &ProxyRouteConfig{
			Prefix: flag.Arg(0),
			Target: flag.Arg(1),
		}})
	} else if flag.NArg() != 0 {
		usage(false, 1)
	}

	return cliConfig
}
