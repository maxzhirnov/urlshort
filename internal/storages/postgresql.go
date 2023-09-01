package storages

import (
	"database/sql"

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

func (p Postgresql) Insert(shortURL models.ShortURL) error {
	return nil
}
func (p Postgresql) Get(id string) (*models.ShortURL, bool) {
	return nil, false
}

func (p Postgresql) Ping() error {
	return p.DB.Ping()
}
