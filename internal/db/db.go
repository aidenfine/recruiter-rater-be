package db

import (
	"context"
	"errors"
	"time"

	"github.com/aidenfine/recruiter-rater/internal/config"
	"github.com/jmoiron/sqlx"
)

func Connect(cfg *config.Config) (*sqlx.DB, error) {
	uri := cfg.PostgresURI
	if uri == "" {
		return nil, errors.New("no databse url")
	}
	db, err := sqlx.Open("postgres", uri)
	if err != nil {
		return nil, err
	}
	// timeout after 10 seconds if database is not reachable
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil

}
