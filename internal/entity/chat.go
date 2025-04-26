package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Chat struct {
	ID        string         `gorm:"type:varchar(255);primaryKey"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null"`
	Content   string         `gorm:"type:text;not null"`
	Output    string         `gorm:"type:text;not null"`
	CreatedAt time.Time      `gorm:"timestamp;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"timestamp;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"timestamp;index"`
	User      User           `gorm:"foreignKey:UserID;references:ID"`
}
