package migrate

import (
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Up(logger *zap.SugaredLogger, db *sqlx.DB) (bool, error) {

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	logger.Info("Initialized driver for migrations")
	if err != nil {
		return false, err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return false, err
	}
	logger.Info("Migrations are found, driver is set")
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return true, nil
		}
		return false, err
	}
	logger.Info("Migrations ran successfully")
	return true, nil

}
