package migrate

import (
	"database/sql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

//Up takes existing zap.SugaredLogger instance and sqlx.DB instance
//Utilizes postgres connection
//Looks up for migrations/ folder in the root
//Logs each migration step
func Up(logger *zap.SugaredLogger, db *sql.DB) (bool, error) {

	driver, err := postgres.WithInstance(db, &postgres.Config{})
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
