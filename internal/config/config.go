package config

import (
	"errors"
	"flag"
	"os"
	"strings"
)

type Config struct {
	serverAddr      string
	baseURL         string
	fileStoragePath string
}

type Builder struct {
	config Config
}

func (b *Builder) WithServerAddr(serverAddr string) *Builder {
	b.config.serverAddr = serverAddr
	return b
}

func (b *Builder) WithBaseURL(baseURL string) *Builder {
	b.config.baseURL = baseURL
	return b
}

func (b *Builder) WithFileStoragePath(fileStoragePath string) *Builder {
	b.config.fileStoragePath = fileStoragePath
	return b
}

func NewFromFlags() (Config, error) {
	var serverAddr string
	flag.Func("a", "Provide address to run the service (default: localhost:8080)", func(v string) error {
		s := strings.Split(v, ":")
		if len(s) != 2 {
			return errors.New("wrong -a parameter value, should be localhost:8080 or similar")
		}
		serverAddr = v
		return nil
	})

	var baseURL string
	flag.StringVar(&baseURL, "b", "unset", "Provide domain which will be used for serving shorten URLs")

	var fileStoragePath string
	flag.StringVar(&fileStoragePath, "f", "/tmp/short-url-db.json", "Provide full path to the file where urls data will be saved")

	flag.Parse()

	var builder Builder
	builder.WithServerAddr(serverAddr).
		WithBaseURL(baseURL).
		WithFileStoragePath(fileStoragePath)

	if builder.config.baseURL == "unset" {
		baseURL = "http://" + serverAddr
		builder.WithBaseURL(baseURL)
	}

	if v, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		builder.WithServerAddr(v)
	}

	if v, ok := os.LookupEnv("BASE_URL"); ok {
		builder.WithBaseURL(v)
	}

	if v, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		builder.WithFileStoragePath(v)
	}

	return builder.config, nil
}

func (c Config) ServerAddr() string {
	return c.serverAddr
}

func (c Config) BaseURL() string {
	return c.baseURL
}

func (c Config) FileStoragePath() string {
	return c.fileStoragePath
}
