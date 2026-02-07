package models

import (
	"time"

	"github.com/google/uuid"
)

type Widget struct {
	ID        uuid.UUID              `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	PageID    uuid.UUID              `gorm:"type:uuid;not null" json:"page_id"`
	Type      string                 `json:"type"`
	Position  int                    `json:"position"`
	Config    map[string]interface{} `gorm:"type:jsonb" json:"config,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}
