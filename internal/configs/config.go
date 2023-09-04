package configs

import (
	"flag"
	"os"
)

type logger interface {
	Info(string, ...interface{})
	Error(string, ...interface{})
	Fatal(string, ...interface{})
	Warn(string, ...interface{})
	Debug(string, ...interface{})
}

const (
	serverAddrFlag      = "a"
	baseURLFlag         = "b"
	fileStoragePathFlag = "f"
	postgresConnFlag    = "d"

	defaultServerAddr      = "localhost:8080"
	defaultBaseURL         = "http://" + defaultServerAddr
	defaultFileStoragePath = "/tmp/short-url-db.json"
	defaultPostgresConn    = ""

	serverAddrFlagUsageMessage  = "Provide server address"
	baseURLFlagUsageMessage     = "Provide domain which will be used for serving shorten URLs"
	fileStoragePathUsageMessage = "Provide full path to the file where urls data will be saved"
	postgresConnUsageMessage    = "Provide PostgreSQL DB connection string"
)

type Config struct {
	serverAddr      string
	baseURL         string
	fileStoragePath string
	postgresConn    string
	logger          logger
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

func (b *Builder) WithPostgresConn(postgresConn string) *Builder {
	b.config.postgresConn = postgresConn
	return b
}

func NewFromFlags(logger logger) (*Config, error) {
	var serverAddr string
	flag.StringVar(&serverAddr, serverAddrFlag, defaultServerAddr, serverAddrFlagUsageMessage)

	var baseURL string
	flag.StringVar(&baseURL, baseURLFlag, defaultBaseURL, baseURLFlagUsageMessage)

	var fileStoragePath string
	flag.StringVar(&fileStoragePath, fileStoragePathFlag, defaultFileStoragePath, fileStoragePathUsageMessage)

	var postgresConn string
	flag.StringVar(&postgresConn, postgresConnFlag, defaultPostgresConn, postgresConnUsageMessage)

	flag.Parse()

	var builder Builder
	builder.WithServerAddr(serverAddr).
		WithBaseURL(baseURL).
		WithFileStoragePath(fileStoragePath).
		WithPostgresConn(postgresConn)

	if v, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		logger.Debug("successfully parsed SERVER_ADDRESS from env")
		builder.WithServerAddr(v)
	} else {
		logger.Warn("couldn't parse SERVER_ADDRESS from env")
	}

	if v, ok := os.LookupEnv("BASE_URL"); ok {
		logger.Debug("successfully parsed BASE_URL from env")
		builder.WithBaseURL(v)
	} else {
		logger.Warn("couldn't parse BASE_URL from env")
	}

	if v, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		logger.Debug("successfully parsed FILE_STORAGE_PATH from env")
		builder.WithFileStoragePath(v)
	} else {
		logger.Warn("couldn't parse FILE_STORAGE_PATH from env")
	}

	if v, ok := os.LookupEnv("POSTGRES_CONN"); ok {
		logger.Debug("successfully parsed POSTGRES_CONN from env")
		builder.WithPostgresConn(v)
	} else {
		logger.Warn("couldn't parse POSTGRES_CONN from env")
	}

	cfg := &builder.config
	cfg.logger = logger

	return cfg, nil
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

func (c Config) ShouldSaveToFile() bool {
	return c.fileStoragePath != ""
}

func (c Config) ShouldUsePostgres() bool {
	return c.postgresConn != ""
}

func (c Config) PostgresConn() string {
	return c.postgresConn
}
