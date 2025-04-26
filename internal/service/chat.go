package service

import (
	"context"

	"github.com/CRobinDev/BCCGembira_Nusastra/internal/dto"
	"github.com/CRobinDev/BCCGembira_Nusastra/internal/entity"
	"github.com/CRobinDev/BCCGembira_Nusastra/internal/repository"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/claude"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/helper"
	"github.com/sirupsen/logrus"
)

type IChatService interface {
	CreateChat(ctx context.Context, req dto.ChatRequest) (dto.ChatResponse, error)
	GetChatByID(ctx context.Context, req dto.ChatHistoryRequest) (dto.ChatResponse, error)
}

type chatService struct {
	cr     repository.IChatRepository
	logger *logrus.Logger
	claude claude.IClaude
}

func NewChatService(cr repository.IChatRepository, logger *logrus.Logger, claude claude.IClaude) IChatService {
	return &chatService{
		cr:     cr,
		logger: logger,
		claude: claude,
	}
}

func (cs *chatService) CreateChat(ctx context.Context, req dto.ChatRequest) (dto.ChatResponse, error) {
	promptResp, err := cs.claude.CreateChat(req)
	if err != nil {
		cs.logger.WithFields(map[string]interface{}{
			"error":   err.Error(),
			"user_id": req.UserID.String(),
		}).Error("[chatService.CreateChat] failed to create chat with claude")
		return dto.ChatResponse{}, err
	}
	time := helper.GetCurrentTime()

	chat := &entity.Chat{
		ID:        promptResp.ID,
		UserID:    req.UserID,
		Content:   req.Content,
		Output:    promptResp.Translation,
		CreatedAt: time,
		UpdatedAt: time,
	}

	err = cs.cr.SaveChat(ctx, chat)
	if err != nil {
		cs.logger.WithFields(map[string]interface{}{
			"error":   err.Error(),
			"user_id": req.UserID.String(),
		}).Error("[chatService.CreateChat] failed to save chat to database")
		return dto.ChatResponse{}, err
	}

	cs.logger.WithFields(map[string]interface{}{
		"id": promptResp.ID,
	}).Info("[chatService.CreateChat] chat created successfully")

	return dto.ChatResponse{
		ID:             chat.ID,
		UserID:         chat.UserID,
		Translation:    promptResp.Translation,
		SourceLanguage: promptResp.SourceLanguage,
		Explanation:    promptResp.Explanation,
	}, nil
}

func (cs *chatService) GetChatByID(ctx context.Context, req dto.ChatHistoryRequest) (dto.ChatResponse, error) {
	chat, err := cs.cr.GetChatByID(ctx, req.ID)
	if err != nil {
		return dto.ChatResponse{}, err
	}

	return dto.ChatResponse{
		ID:      chat.ID,
		UserID:  chat.UserID,
		Translation: chat.Content,
	}, nil
}
