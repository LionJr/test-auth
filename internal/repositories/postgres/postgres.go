package postgres

import (
	"fmt"
	"test-auth/internal/config"

	_ "github.com/jackc/pgx/v4/stdlib" // pgx driver
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func NewPostgresDB(cfg config.Postgres) (*sqlx.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.DBName,
		cfg.Password,
	)

	db, err := sqlx.Connect("pgx", dataSourceName)
	if err != nil {
		return nil, errors.Wrap(err, "sqlx.Connect")
	}

	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "db.Ping")
	}

	return db, nil
}
