package dto

import "github.com/google/uuid"

type ChatRequest struct {
	UserID  uuid.UUID `json:"user_id"`
	Content string    `json:"content" validate:"required"`
	Type    string    `json:"type" validate:"required,oneof=text image"`
}

type ChatResponse struct {
	ID             string    `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	SourceLanguage string    `json:"source_language,omitempty"`
	Translation    string    `json:"translation,omitempty"`
	Explanation    string    `json:"explanation,omitempty"`
}

type ChatHistoryRequest struct {
	ID string `json:"id"`
}
