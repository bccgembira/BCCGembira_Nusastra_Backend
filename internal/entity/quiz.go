package entity

import "time"

type Quiz struct {
	ID        int64      `gorm:"type:uint;primaryKey;unique"`
	Title     string     `gorm:"type:varchar(255);default:null"`
	Language  string     `gorm:"type:varchar(255);default:null"`
	Type      string     `gorm:"type:varchar(255);default:null"`
	Points    int64      `gorm:"type:uint;default:0"`
	CreatedAt time.Time  `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt time.Time  `gorm:"type:timestamp;autoUpdateTime"`
	Questions []Question `gorm:"foreignKey:QuizID;references:ID"` 
}

type QuizAttempt struct {
	ID        int64     `gorm:"type:uint;primaryKey"`
	UserID    string    `gorm:"type:varchar(255);not null"`
	QuizID    int64     `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt time.Time `gorm:"type:timestamp;autoUpdateTime"`
	User      User      `gorm:"foreignKey:UserID;references:ID"`
	Quiz      Quiz      `gorm:"foreignKey:QuizID;references:ID"`
}
