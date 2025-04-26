package entity

import "time"

type Question struct {
	ID      int64            `gorm:"type:uint;primaryKey;autoIncrement"`
	QuizID  int64            `gorm:"type:uint;not null"`
	Title   string           `gorm:"type:varchar(255);not null"`
	Choices []QuestionChoice `gorm:"foreignKey:QuestionID;references:ID"`
	Quiz    Quiz             `gorm:"foreignKey:QuizID;references:ID"`
}

type QuestionChoice struct {
	ID          int64     `gorm:"type:uint;primaryKey;autoIncrement"`
	QuestionID  int64     `gorm:"type:uint;not null"`
	Description string    `gorm:"type:varchar(255);not null"`
	IsCorrect   bool      `gorm:"type:boolean;default:false"`
	CreatedAt   time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"type:timestamp;autoUpdateTime"`
	Question    Question  `gorm:"foreignKey:QuestionID;references:ID"`
}
