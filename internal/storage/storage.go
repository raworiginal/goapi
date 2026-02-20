// Package storage is used to init and update the db
package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/raworiginal/goapi/internal/project"
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
func CreateProject(p *project.Project) error {
	result := DB.Create(p)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetProject retrieves a project by name
func GetProject(name string) (*project.Project, error) {
	var p project.Project
	if err := DB.Where("name = ?", name).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// ListProjects retrieves all projects
func ListProjects() ([]*project.Project, error) {
	var projects []*project.Project
	if err := DB.Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

// DeleteProject removes a project by name
func DeleteProject(name string) error {
	result := DB.Where("name = ?", name).Delete(&project.Project{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("project '%s' not found", name)
	}
	return nil
}
