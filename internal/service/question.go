package service

import (
	"context"

	"github.com/CRobinDev/BCCGembira_Nusastra/internal/dto"
	"github.com/CRobinDev/BCCGembira_Nusastra/internal/repository"
	"github.com/sirupsen/logrus"
)

type IQuestionService interface {
	GetAllQuestionsByQuizID(ctx context.Context, id uint64) ([]dto.QuestionResponse, error)
}

type questionService struct {
	cr     repository.IQuestionRepository
	logger *logrus.Logger
}

func NewQuestionService(cr repository.IQuestionRepository, logger *logrus.Logger) IQuestionService {
	return &questionService{
		cr:     cr,
		logger: logger,
	}
}

func (qs *questionService) GetAllQuestionsByQuizID(ctx context.Context, id uint64) ([]dto.QuestionResponse, error) {
	questions, err := qs.cr.GetAllQuestionsByQuizID(ctx, id)
	if err != nil {
		qs.logger.WithFields(map[string]interface{}{
		"error" : err.Error(),	
	}).Error("[questionService.GetAllQuestionByQuizID] Failde to get questions by quiz id")
		return nil, err
	}

	var questionResponses []dto.QuestionResponse
	for _, question := range questions {
		questionResponses = append(questionResponses, dto.QuestionResponse{
			ID:     question.ID,
			QuizID: question.QuizID,
			Title:  question.Title,
		})
	}

	return questionResponses, nil
}