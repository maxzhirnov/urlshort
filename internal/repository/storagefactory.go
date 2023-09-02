package repository

import (
	"github.com/maxzhirnov/urlshort/internal/configs"
	"github.com/maxzhirnov/urlshort/internal/storages"
)

// NewStorage is a factory function which creating a instance of a storage object of
// significant type and returns it as an Storage interface
func NewStorage(config configs.Config) (Storage, error) {
	switch {
	case config.ShouldUsePostgres():
		postgres, err := storages.NewPostgresql(config.PostgresConn())
		if err != nil {
			return nil, err
		}
		return postgres, nil
	case config.ShouldSaveToFile():
		memStorage := storages.NewMemoryStorage()
		fileStorage, err := storages.NewFileStorage(config.FileStoragePath())
		if err != nil {
			return nil, err
		}
		combined := storages.NewCombinedStorage(memStorage, fileStorage)
		return combined, nil
	default:
		return storages.NewMemoryStorage(), nil
	}
}
