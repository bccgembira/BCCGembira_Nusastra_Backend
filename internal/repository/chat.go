package repository

import (
	"context"
	"errors"

	"github.com/CRobinDev/BCCGembira_Nusastra/internal/entity"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/response"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IChatRepository interface {
	SaveChat(ctx context.Context, chat *entity.Chat) error
	GetChatByID(ctx context.Context, id string) (entity.Chat, error)
	// DeleteOldHistory(userID string) error
}

type chatRepository struct {
	db *gorm.DB
	logger *logrus.Logger
}

func NewChatRepository(db *gorm.DB, logger *logrus.Logger) IChatRepository {
	return &chatRepository{
		db: db,
		logger: logger,
	}
}

func (cr *chatRepository) SaveChat(ctx context.Context, chat *entity.Chat) error {
	err := cr.db.WithContext(ctx).Create(chat).Error
	if err != nil {
		cr.logger.Errorf("failed to save chat: %v", err)
		return &response.ErrSaveChatFailed
	}

	return nil
}

func (cr *chatRepository) GetChatByID(ctx context.Context, id string) (entity.Chat, error) {
	var chat entity.Chat
	if err := cr.db.WithContext(ctx).First(&chat, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cr.logger.Warnf("user with id %s not found", id)
			return entity.Chat{}, &response.ErrUserNotFound
		}
		cr.logger.Errorf("failed to find user by id: %v", err)
		return entity.Chat{}, &response.ErrConnectDatabase
	}

	return chat, nil
}