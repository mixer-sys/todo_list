package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID          uint           `gorm:"primaryKey"`
	Name        string         `gorm:"not null"`
	Description string         `gorm:"type:text"`
	Status      string         `gorm:"not null;default:'pending'"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index;default:null"`
}
