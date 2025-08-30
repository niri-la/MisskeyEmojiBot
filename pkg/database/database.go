package database

import (
	"MisskeyEmojiBot/pkg/entity"
	"MisskeyEmojiBot/pkg/errors"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	*gorm.DB
}

func NewDatabase(dbPath string) (*DB, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := ensureDir(dir); err != nil {
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Suppress GORM logs
	})
	if err != nil {
		return nil, errors.FileOperation("failed to open database", err)
	}

	return &DB{db}, nil
}

func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (db *DB) Migrate() error {
	err := db.AutoMigrate(&entity.Emoji{})
	if err != nil {
		return errors.FileOperation("failed to migrate database", err)
	}

	// Add indexes manually if needed
	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_emojis_is_finish ON emojis(is_finish)").Error
	if err != nil {
		return errors.FileOperation("failed to create indexes", err)
	}

	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_emojis_start_at ON emojis(start_at)").Error
	if err != nil {
		return errors.FileOperation("failed to create indexes", err)
	}

	return nil
}

func ensureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}