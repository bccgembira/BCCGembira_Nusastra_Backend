package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;"`
	DisplayName string         `gorm:"type:varchar(255);not null;unique"`
	Email       string         `gorm:"type:varchar(255);unique;not null;"`
	Password    string         `gorm:"type:varchar(255);not null;"`
	Image       string         `gorm:"type:varchar(255);default:null"`
	Point       int64          `gorm:"type:uint;default:0"`
	CreatedAt   time.Time      `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"type:timestamp;index"`
}

type UserAnswer struct {
	ID            int64       `gorm:"type:uint;primaryKey"`
	UserID        uuid.UUID   `gorm:"type:uuid;not null"`
	QuizAttemptID int64       `gorm:"type:uint;not null"`
	Answer        string      `gorm:"type:varchar(255);not null"`
	Score         int64       `gorm:"type:uint;default:0"`
	CreatedAt     time.Time   `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt     time.Time   `gorm:"type:timestamp;autoUpdateTime"`
	QuizAttempt   QuizAttempt `gorm:"foreignKey:QuizAttemptID;references:ID"`
	User          User        `gorm:"foreignKey:UserID;references:ID"`
}
