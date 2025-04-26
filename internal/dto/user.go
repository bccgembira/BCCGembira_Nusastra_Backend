package dto

import (
	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=24"`
}

type RegisterRequest struct {
	DisplayName string `json:"display_name" validate:"required,min=5,max=255"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8,max=24"`
	Picture     string `json:"picture,omitempty"`
}

type TokenLoginRequest struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

type UpdateRequest struct {
	ID          uuid.UUID `json:"id"`
	DisplayName string    `json:"display_name" validate:"min=5,max=255"`
	NewPassword string    `json:"new_password,omitempty"`
}

type DeleteRequest struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

type NotificationRequest struct {
	ID          uuid.UUID `json:"-"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	Feature     string    `json:"feature" validate:"required,oneof=Blog Portofolio"`
	Link        string    `json:"link" validate:"required"`
}

type LoginResponse struct {
	ID          uuid.UUID `json:"userID"`
	DisplayName string    `json:"display_name"`
	Token       string    `json:"token"`
}

type TokenLoginResponse struct {
	ID          uuid.UUID `json:"userID"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
}
