package storage

import (
	"fmt"

	"github.com/raworiginal/goapi/internal/route"
)

// CreateRoute Adds new route to project in database
func CreateRoute(r *route.Route) error {
	if r.Name == "" {
		return fmt.Errorf("route name cannot be empty")
	}
	result := DB.Create(r)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// ListRoutesByProject retreieves all routes for a project
func ListRoutesByProject(projectID uint) ([]*route.Route, error) {
	var routes []*route.Route
	result := DB.Where("project_id = ?", projectID).Find(&routes)
	if result.Error != nil {
		return nil, result.Error
	}
	return routes, nil
}

// GetRoute retrieves a route by ID
func GetRoute(id uint) (*route.Route, error) {
	var r route.Route
	if err := DB.Where("id = ?", id).First(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

// GetRouteByName retrieves a route by project ID and name
func GetRouteByName(projectID uint, name string) (*route.Route, error) {
	var r route.Route
	if err := DB.Where("name = ? AND project_id = ?", name, projectID).First(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

// UpdateRoute modifies an existing route
func UpdateRoute(id uint, updates *route.UpdateRouteInput) error {
	result := DB.Model(&route.Route{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no route found with id: %v", id)
	}
	return nil
}

// DeleteRoute removes a route by ID
func DeleteRoute(id uint) error {
	result := DB.Where("id = ?", id).Delete(&route.Route{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("route not found: route %v", id)
	}
	return nil
}
