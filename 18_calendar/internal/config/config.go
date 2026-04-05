package config

import (
	"flag"
	"os"
)

type Config struct {
	Port string
}

func Load() Config {
	port := flag.String("port", "", "HTTP server port")
	flag.Parse()

	if *port != "" {
		return Config{Port: *port}
	}

	if p := os.Getenv("PORT"); p != "" {
		return Config{Port: p}
	}

	return Config{Port: "8080"}
}
