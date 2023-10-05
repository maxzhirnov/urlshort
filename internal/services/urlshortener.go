package services

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/maxzhirnov/urlshort/internal/models"
	"github.com/maxzhirnov/urlshort/internal/repositories"
)

const (
	deletionInterval = 10 * time.Second
	deleteChanCap    = 512
)

var (
	ErrEntityAlreadyExist = errors.New("entity already exist")
)

type logger interface {
	Info(string, ...interface{})
	Error(string, ...interface{})
	Fatal(string, ...interface{})
	Warn(string, ...interface{})
	Debug(string, ...interface{})
}

type repository interface {
	Insert(context.Context, models.ShortURL) (models.ShortURL, error)
	InsertMany(context.Context, []models.ShortURL) ([]models.ShortURL, error)
	GetURLByID(ctx context.Context, id string) (models.ShortURL, error)
	GetURLsByUUID(ctx context.Context, uuid string) ([]models.ShortURL, error)
	TagURLsDeleted([]models.Deletion) error
	Ping() error
}

type idGenerator interface {
	Generate() string
}

type URLShortener struct {
	Repo        repository
	IDGenerator idGenerator
	logger      logger

	// Канал для удаления URL-ов
	deleteChan     chan models.Deletion
	deletionsStack []models.Deletion
}

func NewURLShortener(repo repository, idGenerator idGenerator, logger logger) *URLShortener {
	return &URLShortener{
		Repo:           repo,
		IDGenerator:    idGenerator,
		logger:         logger,
		deleteChan:     make(chan models.Deletion, deleteChanCap),
		deletionsStack: make([]models.Deletion, 0, deleteChanCap),
	}
}

func (us *URLShortener) Create(originalURL, uuid string) (models.ShortURL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	urlShorten := models.ShortURL{
		OriginalURL: originalURL,
		ID:          us.IDGenerator.Generate(),
		UUID:        uuid,
	}
	if originalURL == "" {
		return models.ShortURL{}, errors.New("originalURL shouldn't be empty string")
	}
	insertedURL, err := us.Repo.Insert(ctx, urlShorten)

	if errors.Is(err, repositories.ErrEntityAlreadyExist) {
		return insertedURL, ErrEntityAlreadyExist
	}

	if err != nil {
		return models.ShortURL{}, err
	}

	return insertedURL, nil
}

func (us *URLShortener) Get(id string) (models.ShortURL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if id == "" {
		return models.ShortURL{}, errors.New("id shouldn't be empty string")
	}
	return us.Repo.GetURLByID(ctx, id)
}

func (us *URLShortener) CreateBatch(urls []string, uuid string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if len(urls) == 0 {
		return []string{}, nil
	}

	urlsToInsert := make([]models.ShortURL, len(urls))

	for i, url := range urls {
		urlsToInsert[i] = models.ShortURL{
			OriginalURL: url,
			ID:          us.IDGenerator.Generate(),
			UUID:        uuid,
		}
	}

	shortenURLs, err := us.Repo.InsertMany(ctx, urlsToInsert)
	if err != nil {
		us.logger.Error(err.Error())
		return nil, err
	}

	ids := make([]string, len(urls))
	for i, u := range shortenURLs {
		ids[i] = u.ID
	}
	return ids, nil
}

func (us *URLShortener) GetAllUsersURLs(uuid string) ([]models.ShortURL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return us.Repo.GetURLsByUUID(ctx, uuid)
}

func (us *URLShortener) Delete(ids []string, userID string) {
	go func() {
		for _, id := range ids {
			deletion := models.Deletion{
				UserID: userID,
				URLID:  id,
			}
			us.deleteChan <- deletion
		}
	}()
}

func (us *URLShortener) ProcessLinkDeletion(ctx context.Context) {
	ticker := time.NewTicker(deletionInterval)
	defer ticker.Stop()

	for {
		select {
		case d := <-us.deleteChan:
			us.logger.Debug("adding deletion to delete chan")
			us.deletionsStack = append(us.deletionsStack, d)
		case <-ticker.C:
			if len(us.deletionsStack) == 0 {
				continue
			}
			err := us.Repo.TagURLsDeleted(us.deletionsStack)
			if err != nil {
				us.logger.Error(err.Error())
				continue
			}
			us.deletionsStack = us.deletionsStack[0:]
		case <-ctx.Done():
			us.Stop()
		}
	}
}

func (us *URLShortener) Stop() {
	us.logger.Debug("Shutting down...")
	// Удаляем все оставшиеся в канале deletions
	for len(us.deleteChan) > 0 {
		d := <-us.deleteChan
		us.deletionsStack = append(us.deletionsStack, d)
	}
	if len(us.deletionsStack) > 0 {
		err := us.Repo.TagURLsDeleted(us.deletionsStack)
		if err != nil {
			us.logger.Error(err.Error())
		}
	}
	// и закрываем канал
	close(us.deleteChan)
	os.Exit(0)
}

func (us *URLShortener) Ping() error {
	return us.Repo.Ping()
}
