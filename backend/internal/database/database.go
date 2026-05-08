package database

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fitrianabila2025group/videoxnx/backend/internal/models"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect opens the database. The driver is chosen from the DSN prefix:
//
//	postgres://...   or  postgresql://...   → PostgreSQL
//	sqlite://path    or  file:path           → SQLite (file)
//	(empty)          → SQLite at ./data/videoxnx.db (sane default for Railway)
//
// SQLite is pure-Go (no CGO) so it works in scratch/alpine images and Railway
// without any add-on database. Use Postgres in production for higher concurrency.
func Connect(dsn string) (*gorm.DB, error) {
	// Silence the noisy "record not found" warnings emitted on every upsert
	// lookup, but keep real SQL errors visible.
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Error,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
	cfg := &gorm.Config{Logger: gormLogger}

	switch {
	case strings.HasPrefix(dsn, "postgres://"), strings.HasPrefix(dsn, "postgresql://"):
		return gorm.Open(postgres.Open(dsn), cfg)
	default:
		path := strings.TrimPrefix(dsn, "sqlite://")
		path = strings.TrimPrefix(path, "file:")
		if path == "" {
			path = "data/videoxnx.db"
		}
		if dir := filepath.Dir(path); dir != "" && dir != "." {
			_ = os.MkdirAll(dir, 0o755)
		}
		return gorm.Open(sqlite.Open(path+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)"), cfg)
	}
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Tag{},
		&models.Post{},
		&models.ScrapeLog{},
		&models.Report{},
		&models.Setting{},
	)
}
