package storage

import (
	"fmt"

	"github.com/raworiginal/goapi/internal/project"
)

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
