package handler

import (
	"github.com/CRobinDev/BCCGembira_Nusastra/internal/dto"
	"github.com/CRobinDev/BCCGembira_Nusastra/internal/service"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/jwt"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/response"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type IConnectionHandler interface {
	CreateConnection() fiber.Handler
	DeleteConnection() fiber.Handler
	GetAllConnection() fiber.Handler
}

type connectionHandler struct {
	cs  service.IConnectionService
	val validator.Validator
}

func NewConnectionHandler(cs service.IConnectionService, val validator.Validator) IConnectionHandler {
	return &connectionHandler{
		cs:  cs,
		val: val,
	}
}

func (ch *connectionHandler) CreateConnection() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := dto.ConnectionRequest{}
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

		err = ch.cs.CreateConnection(c.Context(), req)
		if err != nil {
			return err
		}

		return response.Success(c, "connections", nil)
	}
}

func (ch *connectionHandler) DeleteConnection() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := dto.ConnectionDeletionRequest{}
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

		err = ch.cs.DeleteConnection(c.Context(), req)
		if err != nil {
			return err
		}

		return response.Success(c, "users", nil)
	}
}

func (ch *connectionHandler) GetAllConnection() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := dto.GetConnectionRequest{}
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

		resp, err := ch.cs.GetAllConnection(c.Context(), req)
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "success",
			"data":    resp,
		})
	}
}
