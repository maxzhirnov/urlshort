package config

import (
	"errors"
	"flag"
	"github.com/maxzhirnov/urlshort/internal/app"
	"strings"
)

type Protocol string

const (
	HTTP  Protocol = "http://"
	HTTPS          = "https://"
)

type Config struct {
	ServerAddr          string
	RedirectHost        string
	RedirectURLProtocol Protocol
}

func NewDefaultConfig() *Config {
	return &Config{
		ServerAddr:          "localhost:8080",
		RedirectHost:        "localhost:8080",
		RedirectURLProtocol: HTTPS,
	}
}

func (c *Config) ParseFlags() {
	flag.Func("a", "Provide address to run the service (e.g. localhost:8080", func(v string) error {
		s := strings.Split(v, ":")
		if len(s) != 2 {
			return errors.New("wrong -a parameter value, should be localhost:8080 or similar")
		}
		c.ServerAddr = v
		return nil
	})

	flag.Func("b", "Provide domain which will be used for serving shorten URLs", func(v string) error {
		u, ok := app.CheckURL(v)
		if !ok {
			return errors.New("-b value should be valid domain")
		}
		c.RedirectHost = u
		return nil
	})

	flag.Parse()
}
