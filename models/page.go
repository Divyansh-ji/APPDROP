package models

import (
	"time"

	"github.com/google/uuid"
)

type page struct {
	ID        uuid.UUID `gorm: "type: uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name      string    `json:"name"`
	Route     string    `json:"route"`
	IsHome    bool      `json:"is_home"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
