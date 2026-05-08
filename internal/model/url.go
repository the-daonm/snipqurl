package model

import "time"

type URL struct {
	ID          int64     `db:"id"`
	OriginalURL string    `db:"original_url"`
	ShortCode   string    `db:"short_code"`
	Clicks      int       `db:"clicks"`
	CreatedAt   time.Time `db:"created_at"`
}
