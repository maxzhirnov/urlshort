package config

import (
	"errors"
	"flag"
	"os"
	"strings"
)

type Protocol string

type Config struct {
	ServerAddr string
	BaseURL    string
}

func NewDefaultConfig() *Config {
	return &Config{
		ServerAddr: "localhost:8080",
		BaseURL:    "http://localhost:8080",
	}
}

func (c *Config) Parse() {
	flag.Func("a", "Provide address to run the service (default: localhost:8080)", func(v string) error {
		s := strings.Split(v, ":")
		if len(s) != 2 {
			return errors.New("wrong -a parameter value, should be localhost:8080 or similar")
		}
		c.ServerAddr = v
		return nil
	})

	flag.StringVar(&c.BaseURL, "b", "unset", "Provide domain which will be used for serving shorten URLs")

	flag.Parse()

	// Если запускаемся только с флагом -a то устанавливаем RedirectHost как http://{ServerAddr}
	// Потому что так кажется логичнее
	if c.BaseURL == "unset" {
		c.BaseURL = "http://" + c.ServerAddr
	}

	if s, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		c.ServerAddr = s
	}

	if s, ok := os.LookupEnv("BASE_URL"); ok {
		c.BaseURL = s
	}
}
