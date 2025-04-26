package service

import (
	"context"

	"github.com/CRobinDev/BCCGembira_Nusastra/internal/dto"
	"github.com/CRobinDev/BCCGembira_Nusastra/internal/entity"
	"github.com/CRobinDev/BCCGembira_Nusastra/internal/repository"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/helper"
	"github.com/sirupsen/logrus"
)

type IConnectionService interface {
	CreateConnection(ctx context.Context, req dto.ConnectionRequest) error
	DeleteConnection(ctx context.Context, req dto.ConnectionDeletionRequest) error
	GetAllConnection(ctx context.Context, req dto.GetConnectionRequest) ([]dto.GetConnectionResponse, error)
}

type connectionService struct {
	cr     repository.IConnectionRepository
	ur     repository.IUserRepository
	logger *logrus.Logger
}

func NewConnectionService(cr repository.IConnectionRepository, ur repository.IUserRepository, logger *logrus.Logger) IConnectionService {
	return &connectionService{
		cr:     cr,
		ur:     ur,
		logger: logger,
	}
}

func (cs *connectionService) CreateConnection(ctx context.Context, req dto.ConnectionRequest) error {
	time := helper.GetCurrentTime()

	connection := entity.Connection{
		UserID:    req.UserID,
		FriendID:  req.FriendID,
		UpdatedAt: time,
	}

	err := cs.cr.CreateConnection(ctx, &connection)
	if err != nil {
		cs.logger.WithFields(map[string]interface{}{
			"error":   err.Error(),
			"user_id": req.UserID.String(),
		}).Error("[connectionService.CreateConnection] failed to create connection")
		return err
	}

	cs.logger.WithFields(map[string]interface{}{
		"user_id":   req.UserID.String(),
		"friend_id": req.FriendID.String(),
	}).Info("[connectionService.CreateConnection] connection created successfully")
	return nil
}

func (cs *connectionService) DeleteConnection(ctx context.Context, req dto.ConnectionDeletionRequest) error {
	connection := entity.Connection{
		UserID:   req.UserID,
		FriendID: req.FriendID,
	}

	err := cs.cr.DeleteConnection(ctx, &connection)
	if err != nil {
		cs.logger.WithFields(map[string]interface{}{
			"error":   err.Error(),
			"user_id": req.UserID.String(),
		}).Error("[connectionService.DeleteConnection] failed to delete connection")

		return err
	}

	cs.logger.WithFields(map[string]interface{}{
		"user_id":   req.UserID.String(),
		"friend_id": req.FriendID.String(),
	}).Info("[connectionService.DeleteConnection] connection deleted successfully")
	return nil
}

func (cs *connectionService) GetAllConnection(ctx context.Context, req dto.GetConnectionRequest) ([]dto.GetConnectionResponse, error) {
	connections, err := cs.cr.GetAllConnection(ctx, req.UserID)
	if err != nil {
		cs.logger.WithFields(map[string]interface{}{
			"error":   err.Error(),
			"user_id": req.UserID.String(),
		}).Error("[connectionService.GetConnectionByUserID] failed to get connection by user ID")
		return nil, err
	}

	if len(connections) == 0 {
		return []dto.GetConnectionResponse{}, nil
	}

	var resp []dto.GetConnectionResponse
	for _, connection := range connections {
		friend, err := cs.ur.FindByID(ctx, connection.FriendID)
		if err != nil {
			cs.logger.WithContext(ctx).WithFields(logrus.Fields{
				"friend_id": connection.FriendID,
				"user_id":   req.UserID,
			}).Error("Failed to get user data for friend ID: ", err)
			continue
		}

		resp = append(resp, dto.GetConnectionResponse{
			FriendID:    friend.ID,
			DisplayName: friend.DisplayName,
			Email:       friend.Email,
		})
	}

	return resp, nil

}
