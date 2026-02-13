// Package storage is used to init and update the db
package storage

import (
	"os"
	"path/filepath"

	"github.com/raworiginal/go-api-cli/internal/project"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB // Global instance of the SQLite database

// InitDB initializes the SQLite database and creates tables
func InitDB() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dbDir := filepath.Join(homeDir, ".config", "goapi")

	if err := os.MkdirAll(dbDir, 0o755); err != nil {
		return err
	}

	dbPath := filepath.Join(dbDir, "goapi.db")

	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}
	if err := DB.AutoMigrate(&project.Project{}); err != nil {
		return err
	}

	return nil
}

// CreateProject adds a new project to the database
