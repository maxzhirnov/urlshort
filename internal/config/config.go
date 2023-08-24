package config

import (
	"flag"
	"os"
)

const (
	serverAddrFlag      = "a"
	baseUrlFlag         = "b"
	fileStoragePathFlag = "f"

	defaultServerAddr      = "localhost:8080"
	defaultBaseURL         = "http://" + defaultServerAddr
	defaultFileStoragePath = "/tmp/short-url-db.json"

	serverAddrFlagUsageMessage  = "Provide server address"
	baseURLFlagUsageMessage     = "Provide domain which will be used for serving shorten URLs"
	fileStoragePathUsageMessage = "Provide full path to the file where urls data will be saved"
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

func NewFromFlags() (*Config, error) {
	var serverAddr string
	flag.StringVar(&serverAddr, serverAddrFlag, defaultServerAddr, serverAddrFlagUsageMessage)

	var baseURL string
	flag.StringVar(&baseURL, baseUrlFlag, defaultBaseURL, baseURLFlagUsageMessage)

	var fileStoragePath string
	flag.StringVar(&fileStoragePath, fileStoragePathFlag, defaultFileStoragePath, fileStoragePathUsageMessage)

	flag.Parse()

	var builder Builder
	builder.WithServerAddr(serverAddr).
		WithBaseURL(baseURL).
		WithFileStoragePath(fileStoragePath)

	if v, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		builder.WithServerAddr(v)
	}

	if v, ok := os.LookupEnv("BASE_URL"); ok {
		builder.WithBaseURL(v)
	}

	if v, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		builder.WithFileStoragePath(v)
	}

	return &builder.config, nil
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
