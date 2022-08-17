package postgres

import (
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     int64
	Database string
}

const dialect = "pgx"

func Connect(c *Config) (*sqlx.DB, error) {

	connStr := fmt.Sprintf("user=%s host=%s port=%d dbname=%s password=%s sslmode=disable", c.User, c.Host, c.Port, c.Database, c.Password)
	db, err := sqlx.Open(dialect, connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
