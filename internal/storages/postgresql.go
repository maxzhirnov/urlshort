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

	if err := initTables(db); err != nil {
		return nil, err
	}

	return &Postgresql{
		DB: db,
	}, nil
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

func (p Postgresql) Insert(shortURL models.ShortURL) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err := p.DB.ExecContext(ctx, `INSERT INTO short_urls(id, original_url) 
	VALUES ($1, $2)`, shortURL.ID, shortURL.OriginalURL); err != nil {
		return err
	}
	return nil
}

func (p Postgresql) Get(id string) (*models.ShortURL, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	row := p.DB.QueryRowContext(ctx, `SELECT id, original_url FROM short_urls WHERE id=$1`, id)
	shortURL := models.ShortURL{}
	err := row.Scan(&shortURL.ID, &shortURL.OriginalURL)
	if err != nil {
		return nil, false
	}
	return &shortURL, true
}

func (p Postgresql) Ping() error {
	return p.DB.Ping()
}

func (p Postgresql) Close() error {
	return p.DB.Close()
}
