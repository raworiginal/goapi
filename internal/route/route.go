// Package route handles the routes for the database
package route

import "time"

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
	ProjectID   uint       `gorm:"foreignKey" json:"project_id"`
	Method      HTTPMethod `json:"method"`
	Path        string     `json:"path"`
	Description string     `json:"description"`
	DateCreated time.Time  `gorm:"autoCreateTime" json:"date_created"`
}
