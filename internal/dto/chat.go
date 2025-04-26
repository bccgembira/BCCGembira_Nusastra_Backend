package dto

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type ChatRequest struct {
	UserID         uuid.UUID `json:"user_id"`
	Content        string    `json:"content" validate:"required"`
	Type           string    `json:"type"`
	SourceLanguage string    `json:"source_language" validate:"required"`
	TargetLanguage string    `json:"target_language" validate:"required"`
}

type ChatResponse struct {
	ID             string    `json:"id,omitempty"`
	UserID         uuid.UUID `json:"user_id"`
	SourceLanguage string    `json:"source_language,omitempty"`
	Translation    string    `json:"translation,omitempty"`
	Explanation    string    `json:"explanation,omitempty"`
}

type ChatHistoryRequest struct {
	ID string `json:"id"`
}

type ChatImageRequest struct {
	UserID uuid.UUID             `json:"user_id"`
	File   *multipart.FileHeader `json:"file,omitempty"`
	Url    string                `json:"url,omitempty"`
}

type ChatOCRResponse struct {
	ParsedResults []OCRParsedResult `json:"ParsedResults"`
	OCRExitCode   int               `json:"OCRExitCode"`
}

type OCRParsedResult struct {
	ParsedText   string `json:"ParsedText"`
	ErrorMessage string `json:"ErrorMessage"`
}

type OCRRequest struct {
	File     *multipart.FileHeader `json:"file"`
	Filetype string                `json:"filetype"`
}
