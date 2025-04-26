package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"

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
	CreateChatWithOCR(ctx context.Context, req dto.ChatImageRequest) (dto.ChatResponse, error)
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
	req.Type = "text"
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
		ID:          chat.ID,
		UserID:      chat.UserID,
		Translation: chat.Content,
	}, nil
}

func (cs *chatService) CreateChatWithOCR(ctx context.Context, req dto.ChatImageRequest) (dto.ChatResponse, error) {
	src, err := req.File.Open()
	if err != nil {
		cs.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"user_id": req.UserID.String(),
		}).Error("[chatService.CreateChatWithOCR] failed to open file")
		return dto.ChatResponse{}, err
	}
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		cs.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"user_id": req.UserID.String(),
		}).Error("[chatService.CreateChatWithOCR] failed to read file")
		return dto.ChatResponse{}, err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("language", "eng")

	part, err := writer.CreateFormFile("file", req.File.Filename)
	if err != nil {
		cs.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"user_id": req.UserID.String(),
		}).Error("[chatService.CreateChatWithOCR] failed to create form file")
		return dto.ChatResponse{}, err
	}

	_, err = part.Write(fileBytes)
	if err != nil {
		cs.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"user_id": req.UserID.String(),
		}).Error("[chatService.CreateChatWithOCR] failed to write file to form")
		return dto.ChatResponse{}, err
	}

	err = writer.Close()
	if err != nil {
		cs.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"user_id": req.UserID.String(),
		}).Error("[chatService.CreateChatWithOCR] failed to close writer")
		return dto.ChatResponse{}, err
	}

	request, err := http.NewRequest(http.MethodPost, os.Getenv("OCR_URL"), body)
	if err != nil {
		cs.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"user_id": req.UserID.String(),
		}).Error("[chatService.CreateChatWithOCR] failed to create request")
		return dto.ChatResponse{}, err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("apikey", os.Getenv("OCR_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		cs.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"user_id": req.UserID.String(),
		}).Error("[chatService.CreateChatWithOCR] failed to send request")
		return dto.ChatResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		cs.logger.WithFields(logrus.Fields{
			"error":   "failed to process image",
			"user_id": req.UserID.String(),
			"status":  resp.Status,
		}).Error("[chatService.CreateChatWithOCR] failed to process image")
		return dto.ChatResponse{}, err
	}

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		cs.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"user_id": req.UserID.String(),
		}).Error("[chatService.CreateChatWithOCR] failed to read response body")
		return dto.ChatResponse{}, err
	}

	var ocrResponse dto.ChatOCRResponse
	err = json.Unmarshal(bodyResp, &ocrResponse)
	if err != nil {
		cs.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"user_id": req.UserID.String(),
		}).Error("[chatService.CreateChatWithOCR] failed to unmarshal response body")
		return dto.ChatResponse{}, err
	}

	if len(ocrResponse.ParsedResults) == 0 || ocrResponse.ParsedResults[0].ParsedText == "" {
		cs.logger.WithFields(logrus.Fields{
			"error":   "empty OCR result",
			"user_id": req.UserID.String(),
		}).Error("[chatService.CreateChatWithOCR] OCR result is empty")
		return dto.ChatResponse{}, err
	}

	claudeReq := dto.ChatRequest{
		UserID:  req.UserID,
		Content: ocrResponse.ParsedResults[0].ParsedText,
		Type:    "image",
	}

	claudeResp, err := cs.claude.CreateChat(claudeReq)
	if err != nil {
		cs.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"user_id": req.UserID.String(),
		}).Error("[chatService.CreateChatWithOCR] failed to create chat with claude")
		return dto.ChatResponse{}, err
	}

	cs.logger.WithFields(logrus.Fields{
		"id": claudeResp.ID,
	}).Info("[chatService.CreateChatWithOCR] chat created successfully")

	return dto.ChatResponse{
		ID:             claudeResp.ID,
		UserID:         req.UserID,
		Translation:    claudeResp.Translation,
		SourceLanguage: claudeResp.SourceLanguage,
		Explanation:    claudeResp.Explanation,
	}, nil
}
