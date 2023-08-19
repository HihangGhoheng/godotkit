package gdk_postgre

import (
	"database/sql"
	"fmt"
	gdk_helpers "github.com/HihangGhoheng/godotkit/helpers"
	gdk_types "github.com/HihangGhoheng/godotkit/types"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gl "gorm.io/gorm/logger"
)

type Postgres struct {
	Database *gorm.DB
	SQLDb    *sql.DB
}

func (c *Postgres) ShutdownConnection(logger *logrus.Logger) {
	if err := c.SQLDb.Close(); err != nil {
		gdk_helpers.FailOnError(err, "Failed to shutdown connection of postgresql")
	} else {
		logger.Infof("Successfully close postgre connection!")
	}
}

func Connect(opt gdk_types.PostgreSQLConfig, logger *logrus.Logger) *Postgres {
	if logger == nil {
		logger = logrus.New()
	}
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		opt.Host,
		opt.Username,
		opt.Password,
		opt.DatabaseName,
		opt.Port,
		opt.SSLMode,
	)
	if opt.Password == "" {
		dsn = fmt.Sprintf(
			"host=%s user=%s port=%d dbname=%s sslmode=%s",
			opt.Host,
			opt.Username,
			opt.Port,
			opt.DatabaseName,
			opt.SSLMode,
		)
	}
	db, err := gorm.Open(
		postgres.New(postgres.Config{DSN: dsn}),
		&gorm.Config{Logger: gl.Default.LogMode(gl.Info)},
	)
	if err != nil {
		gdk_helpers.FatalOnError(err, "Error connect to PostgreSQL database")
	}
	sqlDb, err := db.DB()
	if err != nil {
		gdk_helpers.FailOnError(err, "Failed connect to the database")
	}
	if err = sqlDb.Ping(); err != nil {
		gdk_helpers.FailOnError(err, "Failed connect to the database")
	}
	return &Postgres{
		Database: db,
		SQLDb:    sqlDb,
	}
}
