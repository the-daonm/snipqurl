package database

import (
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewDB() (*sqlx.DB, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.user, cfg.password, cfg.host, cfg.port, cfg.name)

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("sqlx open: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database ping: %w", err)
	}

	return db, nil
}

type config struct {
	host, port, user, password, name string
}

func loadConfig() (config, error) {
	required := map[string]string{
		"DB_HOST":     os.Getenv("DB_HOST"),
		"DB_PORT":     os.Getenv("DB_PORT"),
		"DB_USER":     os.Getenv("DB_USER"),
		"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
		"DB_NAME":     os.Getenv("DB_NAME"),
	}

	for key, val := range required {
		if val == "" {
			return config{}, fmt.Errorf("required environment variable %s is not set", key)
		}
	}

	return config{
		host:     required["DB_HOST"],
		port:     required["DB_PORT"],
		user:     required["DB_USER"],
		password: required["DB_PASSWORD"],
		name:     required["DB_NAME"],
	}, nil
}
