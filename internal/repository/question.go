package repository

import (
	"context"

	"github.com/CRobinDev/BCCGembira_Nusastra/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IQuestionRepository interface {
	GetAllQuestionsByQuizID(ctx context.Context, id uint64) ([]entity.Question, error)
}

type questionRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewQuestionRepository(db *gorm.DB, logger *logrus.Logger) IQuestionRepository {
	return &questionRepository{
		db: db,
		logger: logger,
	}
}

func (qr *questionRepository) GetAllQuestionsByQuizID(ctx context.Context, id uint64) ([]entity.Question, error) {
	questions := []entity.Question{}
	err := qr.db.WithContext(ctx).Find(&questions).Error
	if err != nil {
		qr.logger.Errorf("failed to get questions by quiz id %d: %v", id, err)
		return nil, err
	}
	return questions, nil
}
