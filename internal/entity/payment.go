package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payment struct {
	OrderID   string         `json:"id" gorm:"type:varchar(255);primaryKey;"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid; not null"`
	Amount    float64        `json:"amount" gorm:"type:numeric(12,3);not null"`
	Status    string         `json:"status" gorm:"type:varchar(255);default:'pending'"`
	Type      string         `json:"type" gorm:"type:varchar(255);default:null"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"type:timestamp;index"`
	User      User           `gorm:"foreignKey:UserID;references:ID"`
}
