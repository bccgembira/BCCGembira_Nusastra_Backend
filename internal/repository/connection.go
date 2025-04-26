package repository

import (
	"context"

	"github.com/CRobinDev/BCCGembira_Nusastra/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IConnectionRepository interface {
	CreateConnection(ctx context.Context, connection *entity.Connection) error
	DeleteConnection(ctx context.Context, connection *entity.Connection) error
	GetAllConnection(ctx context.Context, userID uuid.UUID) ([]entity.Connection, error)
}

type connectionRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewConnectionRepository(db *gorm.DB, logger *logrus.Logger) IConnectionRepository {
	return &connectionRepository{
		db:     db,
		logger: logger,
	}
}

func (cr *connectionRepository) CreateConnection(ctx context.Context, connection *entity.Connection) error {
	err := cr.db.WithContext(ctx).Create(connection).Error
	if err != nil {
		cr.logger.WithContext(ctx).Error("Failed to create connection: ", err)
		return err
	}

	return nil
}

func (cr *connectionRepository) DeleteConnection(ctx context.Context, connection *entity.Connection) error {
	err := cr.db.WithContext(ctx).Delete(connection).Error
	if err != nil {
		cr.logger.WithContext(ctx).Error("Failed to delete connection: ", err)
		return err
	}

	return nil
}

func (cr *connectionRepository) GetAllConnection(ctx context.Context, userID uuid.UUID) ([]entity.Connection, error) {
	var connections []entity.Connection
	err := cr.db.WithContext(ctx).Find(&connections).Error
	if err != nil {
		cr.logger.Error("Failed to get all connections: ", err)
		return nil, err
	}

	return connections, nil
}
