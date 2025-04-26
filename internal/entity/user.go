package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;"`
	DisplayName string         `gorm:"type:varchar(255);not null;"`
	Email       string         `gorm:"type:varchar(255);unique;not null;"`
	Password    string         `gorm:"type:varchar(255);not null;"`
	Image       string         `gorm:"type:varchar(255);default:null"`
	CreatedAt   time.Time      `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"type:timestamp;index"`
}
