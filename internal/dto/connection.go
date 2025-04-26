package dto

import (
	"github.com/google/uuid"
)

type ConnectionRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	FriendID uuid.UUID `json:"friend_id,omitempty"`
}

type ConnectionDeletionRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	FriendID uuid.UUID `json:"friend_id" validate:"required"`
}

type GetConnectionResponse struct {
	FriendID    uuid.UUID `json:"friend_id"`
	DisplayName string    `json:"name"`
	Email       string    `json:"email"`
}

type GetConnectionRequest struct {
	UserID uuid.UUID `json:"user_id"`
}
