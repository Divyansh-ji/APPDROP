package models

import (
	"time"

	"github.com/google/uuid"
)

type Brand struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name          string    `gorm:"not null" json:"name"`
	Logo          string    `json:"logo"`
	OfficeAddress string    `json:"office_address" gorm:"column:office_address"`
	Domain        string    `gorm:"not null" json:"domain"`
	Email         string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash  string    `json:"-"` // Never return password hash in JSON
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (Brand) TableName() string { return "brands" }
