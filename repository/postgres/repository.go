package postgres

import (
	"database/sql"
	"time"

	"github.com/go-url-shortener/shortener"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(dbURL string) (shortener.RedirectRepository, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewPostgresRepository")
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewPostgresRepository")
	}

	repo := &postgresRepository{
		db: db,
	}

	return repo, nil
}

func (r *postgresRepository) Find(code string) (*shortener.Redirect, error) {
	query := "SELECT code, url, created_at FROM redirects WHERE code = $1"
	row := r.db.QueryRow(query, code)
	redirect := &shortener.Redirect{}
	var createdAt time.Time
	err := row.Scan(&redirect.Code, &redirect.URL, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(shortener.ErrRedirectNotFound, "repository.Redirect.Find")
		}
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}

	redirect.CreatedAt = createdAt.Unix()
	return redirect, nil
}

func (r *postgresRepository) Store(redirect *shortener.Redirect) error {
	query := "INSERT INTO redirects (code, url, created_at) VALUES ($1, $2, $3)"
	_, err := r.db.Exec(query, redirect.Code, redirect.URL, time.Unix(redirect.CreatedAt, 0))
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}

	return nil
}
