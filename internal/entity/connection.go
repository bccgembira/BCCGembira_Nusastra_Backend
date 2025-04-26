package entity

import (
	"time"

	"github.com/google/uuid"
)

type Connection struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	FriendID  uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt time.Time `gorm:"type:timestamp;autoUpdateTime"`
}
