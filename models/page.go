package models

import (
	"time"

	"github.com/google/uuid"
)

type Page struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name      string        `json:"name"`
	Route     string        `json:"route"`
	IsHome    bool          `json:"is_home"`
	Widgets   []Widget `gorm:"foreignKey:PageID" json:"widgets,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
