package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"snipqurl/internal/model"

	"github.com/jmoiron/sqlx"
)

type URLRepository interface {
	Save(url *model.URL) error
	FindByShortCode(code string) (*model.URL, error)
	IncrementClick(code string) error
	DeleteExpired() (int64, error)
}

type urlRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) URLRepository {
	return &urlRepository{db: db}
}

func (r *urlRepository) Save(url *model.URL) error {
	query := `
		INSERT INTO urls (original_url, short_code, expires_at)
		VALUES ($1, $2, $3)`

	_, err := r.db.Exec(query, url.OriginalURL, url.ShortCode, url.ExpiresAt)
	if err != nil {
		return fmt.Errorf("failed to insert url: %w", err)
	}
	return nil
}

var ErrNotFound = errors.New("url not found")

func (r *urlRepository) FindByShortCode(code string) (*model.URL, error) {
	var u model.URL
	query := `SELECT * from urls WHERE short_code = $1`

	err := r.db.Get(&u, query, code)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &u, nil
}

func (r *urlRepository) IncrementClick(code string) error {
	query := `UPDATE urls SET clicks = clicks + 1 WHERE short_code = $1`

	result, err := r.db.Exec(query, code)
	if err != nil {
		return fmt.Errorf("failed to increment clicks: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no url found with code: %s", code)
	}

	return nil
}

func (r *urlRepository) DeleteExpired() (int64, error) {
	query := `DELETE FROM urls WHERE expires_at IS NOT NULL AND expires_at < NOW()`

	result, err := r.db.Exec(query)
	if err != nil {
		return 0, fmt.Errorf("failed to delete expired urls: %w", err)
	}

	rows, _ := result.RowsAffected()
	return rows, nil
}
