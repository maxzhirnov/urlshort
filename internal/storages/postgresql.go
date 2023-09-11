package storages

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/maxzhirnov/urlshort/internal/models"
)

type Postgresql struct {
	DB *sql.DB
}

func NewPostgresql(conn string) (*Postgresql, error) {
	db, err := sql.Open("pgx", conn)
	if err != nil {
		return nil, err
	}

	return &Postgresql{
		DB: db,
	}, nil
}

func (s Postgresql) Insert(ctx context.Context, shortURL models.ShortURL) (models.ShortURL, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		tx.Rollback()
		return models.ShortURL{}, err
	}

	stmt, err := tx.PrepareContext(ctx, `
	INSERT INTO short_urls(id, original_url, updated_at) 
	VALUES ($1, $2, NOW()) 
	ON CONFLICT (original_url) DO UPDATE SET updated_at = NOW()
	RETURNING *, (xmax = 0) AS is_inserted;
	`)
	if err != nil {
		return models.ShortURL{}, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, shortURL.ID, shortURL.OriginalURL)
	if row.Err() != nil {
		tx.Rollback()
		return models.ShortURL{}, fmt.Errorf("something went wrong")
	}

	var result models.ShortURL
	var createdAt time.Time
	var isInserted bool
	if err := row.Scan(&result.ID, &result.OriginalURL, &createdAt, &isInserted); err != nil {
		tx.Rollback()
		return models.ShortURL{}, err
	}

	tx.Commit()
	if !isInserted {
		return result, ErrEntityAlreadyExist
	}

	return result, nil
}

func (s Postgresql) InsertMany(ctx context.Context, urls []models.ShortURL) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO short_urls(id, original_url) VALUES ($1, $2) ON CONFLICT DO NOTHING ")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, url := range urls {
		if _, err := stmt.ExecContext(ctx, url.ID, url.OriginalURL); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s Postgresql) Get(ctx context.Context, id string) (models.ShortURL, bool) {
	row := s.DB.QueryRowContext(ctx, `SELECT id, original_url FROM short_urls WHERE id=$1`, id)
	shortURL := models.ShortURL{}
	err := row.Scan(&shortURL.ID, &shortURL.OriginalURL)
	if err != nil {
		return shortURL, false
	}
	return shortURL, true
}

func (s Postgresql) Bootstrap(ctx context.Context) error {
	// Создаем таблицу short_urls
	if err := s.initTables(); err != nil {
		return err
	}

	// Создаем уникальный индекс для original_url
	if err := s.createUniqueOriginalURLIndex(); err != nil {
		return err
	}

	return nil
}

func (s Postgresql) Ping() error {
	return s.DB.Ping()
}

func (s Postgresql) Close() error {
	return s.DB.Close()
}

func (s Postgresql) initTables() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err := s.DB.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS short_urls (
									  id varchar(20) NOT NULL,
									  original_url varchar(450) NOT NULL,
									  updated_at TIMESTAMP DEFAULT NOW(),
									  PRIMARY KEY (id)) ;`); err != nil {
		return err
	}
	return nil
}

func (s Postgresql) createUniqueOriginalURLIndex() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := "CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_short_url ON short_urls (original_url)"

	if _, err := s.DB.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}
