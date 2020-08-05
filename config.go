package main

type Config struct {
	Port int `json:"port"`
}

func defaultConfig() *Config {
	return &Config{
		Port: 8080,
	}
}

func blankConfig() *Config {
	return &Config{
		Port: -1,
	}
}
