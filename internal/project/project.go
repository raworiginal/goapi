// Package project
package project

import "time"

type Project struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"unique" json:"name"`
	BaseURL     string    `gorm:"column:base_url" json:"base_url"`
	DateCreated time.Time `gorm:"autoCreateTime" json:"date_created"`
	Description string    `json:"description"` // optional
}
