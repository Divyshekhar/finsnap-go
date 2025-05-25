package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type Expense struct {
	ID          string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title       string    `gorm:"not null"`
	Amount      float64   `gorm:"not null"`
	Type        string    `gorm:"default:'expense';not null"`
	Date        time.Time `gorm:"not null"`
	Category    string    `gorm:"not null"`
	Description string

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	UserID string `gorm:"type:uuid;not null"`
	User   User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (i *Expense) BeforeSave(tx *gorm.DB) (err error) {
	i.Category = strings.ToLower(i.Category)
	return
}
