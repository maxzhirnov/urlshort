package storages

import (
	"context"
	"database/sql"
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

func (s Postgresql) Insert(ctx context.Context, shortURL models.ShortURL) error {
	if _, err := s.DB.ExecContext(ctx, `INSERT INTO short_urls(id, original_url) 
	VALUES ($1, $2)`, shortURL.ID, shortURL.OriginalURL); err != nil {
		return err
	}
	return nil
}

func (s Postgresql) InsertMany(ctx context.Context, urls []models.ShortURL) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO short_urls(id, original_url) VALUES ($1, $2)")
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

func (s Postgresql) Get(ctx context.Context, id string) (*models.ShortURL, bool) {
	row := s.DB.QueryRowContext(ctx, `SELECT id, original_url FROM short_urls WHERE id=$1`, id)
	shortURL := models.ShortURL{}
	err := row.Scan(&shortURL.ID, &shortURL.OriginalURL)
	if err != nil {
		return nil, false
	}
	return &shortURL, true
}

func (s Postgresql) Bootstrap(ctx context.Context) error {
	return initTables(s.DB)
}

func (s Postgresql) Ping() error {
	return s.DB.Ping()
}

func (s Postgresql) Close() error {
	return s.DB.Close()
}

func initTables(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS short_urls (
									  id varchar(20) NOT NULL,
									  original_url varchar(450) NOT NULL,
									  PRIMARY KEY (id));`); err != nil {
		return err
	}
	return nil
}
