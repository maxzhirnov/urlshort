package config

import (
	"errors"
	"flag"
	"os"
	"strings"
)

const (
	defaultAddress = "localhost:8080"
	defaultBaseURL = "http://localhost:8080"
)

type Protocol string

type Config struct {
	ServerAddr string
	BaseURL    string
}

func NewDefaultConfig() *Config {
	return &Config{
		ServerAddr: defaultAddress,
		BaseURL:    defaultBaseURL,
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

	if v, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		c.ServerAddr = v
	}

	if v, ok := os.LookupEnv("BASE_URL"); ok {
		c.BaseURL = v
	}
}