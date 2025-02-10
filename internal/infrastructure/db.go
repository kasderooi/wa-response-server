package infrastructure

import (
	"database/sql"
	"time"
	"whatsapp-bot/internal/models"

	"gorm.io/driver/mysql"
	gormLogger "gorm.io/gorm/logger"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type Dialect int

const (
	PostgresSQL Dialect = iota
	SQLServer
	MySQL
	SQLite
)

type Database struct {
	MaxConnLifetime time.Duration
	MaxIdleConns    int
	MaxOpenConns    int
	DataSourceName  string
	Dialect         Dialect
}

func InitDb(dbConfig Database) (*gorm.DB, *sql.DB, error) {
	db, err := openGormDB(dbConfig.Dialect, dbConfig.DataSourceName)
	if err != nil {
		return nil, nil, err
	}
	sqlDb, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	sqlDb.SetMaxOpenConns(dbConfig.MaxOpenConns)
	sqlDb.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDb.SetConnMaxLifetime(dbConfig.MaxConnLifetime)

	return db, sqlDb, nil
}

func MigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(&models.Session{})
}

func openGormDB(dialect Dialect, dsn string) (*gorm.DB, error) {
	cfg := &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Error),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	switch dialect {
	case SQLServer:
		return gorm.Open(sqlserver.Open(dsn), cfg)
	case SQLite:
		return gorm.Open(sqlite.Open(dsn), cfg)
	case MySQL:
		return gorm.Open(mysql.Open(dsn), cfg)
	default:
		return gorm.Open(postgres.Open(dsn), cfg)
	}
}
