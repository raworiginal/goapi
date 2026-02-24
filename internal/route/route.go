// Package route handles the routes for the database
package route

import (
	"fmt"
	"strings"
	"time"
)

type HTTPMethod string

const (
	GET    HTTPMethod = "GET"
	POST   HTTPMethod = "POST"
	PUT    HTTPMethod = "PUT"
	DELETE HTTPMethod = "DELETE"
	PATCH  HTTPMethod = "PATCH"
)

type Route struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	ProjectID   uint       `gorm:"foreignKey; uniqueIndex:idx_project_route_name" json:"project_id"`
	Name        string     `gorm:"uniqueIndex:idx_project_route_name" json:"name"`
	Method      HTTPMethod `json:"method"`
	Path        string     `json:"path"`
	Description string     `json:"description"`
	DateCreated time.Time  `gorm:"autoCreateTime" json:"date_created"`
}

type UpdateRouteInput struct {
	Name        *string     `json:"name,omitempty"`
	Method      *HTTPMethod `json:"method,omitempty"`
	Path        *string     `json:"path,omitempty"`
	Description *string     `json:"description,omitempty"`
}

func ParseHTTPMethod(s string) (HTTPMethod, error) {
	switch strings.ToUpper(s) {
	case "GET":
		return GET, nil
	case "POST":
		return POST, nil
	case "PUT":
		return PUT, nil
	case "DELETE":
		return DELETE, nil
	case "PATCH":
		return PATCH, nil
	default:
		return "", fmt.Errorf("invalid HTTP method: %s. Valid methods are: GET, POST, PUT, PATCH, DELETE", s)
	}
}
