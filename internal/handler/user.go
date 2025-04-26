package handler

import (
	"github.com/CRobinDev/BCCGembira_Nusastra/internal/dto"
	"github.com/CRobinDev/BCCGembira_Nusastra/internal/service"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/jwt"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/response"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type IUserHandler interface {
	Register() fiber.Handler
	Login() fiber.Handler
	GetUser() fiber.Handler
	Update() fiber.Handler
	Delete() fiber.Handler
	Notification() fiber.Handler
	UploadProfileImage() fiber.Handler
}

type userHandler struct {
	us  service.IUserService
	val validator.Validator
}

func NewUserHandler(us service.IUserService, val validator.Validator) IUserHandler {
	return &userHandler{
		us:  us,
		val: val,
	}
}

func (uh *userHandler) Register() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := dto.RegisterRequest{}
		err := c.BodyParser(&req)
		if err != nil {
			return err
		}

		valErr := uh.val.Validate(req)
		if valErr != nil {
			return valErr
		}

		err = uh.us.Register(c.Context(), req)
		if err != nil {
			return err
		}

		return response.Success(c, "users", nil)
	}
}

func (uh *userHandler) Login() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req dto.LoginRequest
		err := c.BodyParser(&req)
		if err != nil {
			return err
		}

		valErr := uh.val.Validate(req)
		if valErr != nil {
			return valErr
		}

		resp, err := uh.us.Login(c.Context(), req)
		if err != nil {
			return err
		}

		return response.Success(c, "users", resp)
	}
}

func (uh *userHandler) Delete() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req dto.DeleteRequest
		userID, err := jwt.GetUser(c)
		if err != nil {
			return err
		}

		req.ID = userID
		valErr := uh.val.Validate(req)
		if valErr != nil {
			return valErr
		}

		err = uh.us.Delete(c.Context(), req)
		if err != nil {
			return err
		}

		return response.Success(c, "users", nil)
	}
}

func (uh *userHandler) GetUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := jwt.GetUser(c)
		if err != nil {
			return err
		}

		req := dto.TokenLoginRequest{
			ID: userID,
		}

		resp, err := uh.us.GetUser(c.Context(), req)
		if err != nil {
			return err
		}

		return response.Success(c, "users", resp)
	}
}

func (uh *userHandler) Update() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req dto.UpdateRequest
		err := c.BodyParser(&req)
		if err != nil {
			return err
		}

		userID, err := jwt.GetUser(c)
		if err != nil {
			return err
		}

		req.ID = userID
		valErr := uh.val.Validate(req)
		if valErr != nil {
			return valErr
		}

		err = uh.us.Update(c.Context(), req)
		if err != nil {
			return err
		}

		return response.Success(c, "users", nil)

	}
}

func (uh *userHandler) Notification() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req dto.NotificationRequest
		err := c.BodyParser(&req)
		if err != nil {
			return err
		}
		userID, err := jwt.GetUser(c)
		if err != nil {
			return err
		}

		req.ID = userID
		valErr := uh.val.Validate(req)
		if valErr != nil {
			return valErr
		}

		err = uh.us.SendNotification(c.Context(), req)
		if err != nil {
			return err
		}

		return response.Success(c, "users", nil)
	}
}

func (ur *userHandler) UploadProfileImage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}

		userID, err := jwt.GetUser(c)
		if err != nil {
			return err
		}

		valErr := ur.val.Validate(file)
		if valErr != nil {
			return valErr
		}

		url, err := ur.us.UploadProfileImage(c.Context(), userID, file)
		if err != nil {
			return err
		}

		return response.Success(c, "users", url)
	}
}
