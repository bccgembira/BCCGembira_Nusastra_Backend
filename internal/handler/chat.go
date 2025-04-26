package handler

import (
	"github.com/CRobinDev/BCCGembira_Nusastra/internal/dto"
	"github.com/CRobinDev/BCCGembira_Nusastra/internal/service"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/jwt"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/response"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type IChatHandler interface {
	CreateChat() fiber.Handler
	GetChatByID() fiber.Handler
	CreateChatWithOCR() fiber.Handler
}

type chatHandler struct {
	chatService service.IChatService
	val         validator.Validator
}

func NewChatHandler(chatService service.IChatService, val validator.Validator) IChatHandler {
	return &chatHandler{
		chatService: chatService,
		val:         val,
	}
}

func (ch *chatHandler) CreateChat() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := dto.ChatRequest{}
		err := c.BodyParser(&req)
		if err != nil {
			return err
		}

		userID, err := jwt.GetUser(c)
		if err != nil {
			return err
		}

		req.UserID = userID
		valErr := ch.val.Validate(req)
		if valErr != nil {
			return valErr
		}

		resp, err := ch.chatService.CreateChat(c.Context(), req)
		if err != nil {
			return err
		}

		return response.Success(c, "chats", resp)
	}
}

func (ch *chatHandler) GetChatByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := dto.ChatHistoryRequest{}
		err := c.BodyParser(&req)
		if err != nil {
			return err
		}

		valErr := ch.val.Validate(req)
		if valErr != nil {
			return valErr
		}

		resp, err := ch.chatService.GetChatByID(c.Context(), req)
		if err != nil {
			return err
		}

		return response.Success(c, "chats", resp)
	}
}

func (ch *chatHandler) CreateChatWithOCR() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := dto.ChatImageRequest{}
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}

		userID, err := jwt.GetUser(c)
		if err != nil {
			return err
		}

		req.UserID = userID
		valErr := ch.val.Validate(file)
		if valErr != nil {
			return valErr
		}

		req.File = file
		resp, err := ch.chatService.CreateChatWithOCR(c.Context(), req)
		if err != nil {
			return err
		}

		return response.Success(c, "chats", resp)
	}
}
