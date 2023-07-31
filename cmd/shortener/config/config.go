package config

import (
	"errors"
	"flag"
	"strings"
)

type Protocol string

type Config struct {
	ServerAddr   string
	RedirectHost string
}

func NewDefaultConfig() *Config {
	return &Config{
		ServerAddr:   "localhost:8080",
		RedirectHost: "http://localhost:8080",
	}
}

func (c *Config) ParseFlags() {
	flag.Func("a", "Provide address to run the service (default: localhost:8080)", func(v string) error {
		s := strings.Split(v, ":")
		if len(s) != 2 {
			return errors.New("wrong -a parameter value, should be localhost:8080 or similar")
		}
		c.ServerAddr = v
		return nil
	})

	flag.Func("b", "Provide domain which will be used for serving shorten URLs (default: http://localhost:8080)", func(v string) error {
		c.RedirectHost = v
		return nil
	})

	flag.Parse()
}
