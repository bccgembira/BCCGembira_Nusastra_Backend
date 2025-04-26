package config

import (
	"encoding/json"
	"errors"

	"github.com/CRobinDev/BCCGembira_Nusastra/internal/middleware"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/response"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/log"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/sirupsen/logrus"
)

func NewFiber() *fiber.App {
	app := fiber.New(fiber.Config{
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		ErrorHandler: newErrorHandler(),
	})

	app.Use(func(c *fiber.Ctx) error {
		if c.Method() == fiber.MethodOptions {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.Next()
	})

	logger := log.NewLogger()
	app.Use(middleware.RateLimiter())
	app.Use(healthcheck.New())
	app.Use(middleware.Helmet())
	app.Use(middleware.Logger(logger))
	app.Use(middleware.Cors())

	return app
}

func newErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {

		var apiError *response.Errors
		if errors.As(err, &apiError) {
			return c.Status(apiError.Code).JSON(fiber.Map{
				"error":   err,
				"message": apiError.Error(),
			})
		}

		var fiberError *fiber.Error
		if errors.As(err, &fiberError) {
			return c.Status(fiberError.Code).JSON(fiber.Map{
				"message": utils.StatusMessage(fiberError.Code),
				"error":   err,
			})
		}

		var validationError validator.ValidationErrors
		if errors.As(err, &validationError) {
			validationDetails := fiber.Map{}
			for field, msg := range validationError {
				validationDetails[field] = msg
			}

			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"errors": validationDetails,
			})
		}

		logrus.WithFields(map[string]interface{}{
			"error":   err.Error(),
			"message": "Congratulations! You've encountered an unhandled error.",
		}).Error("[Error Handler} Unhandled error")

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": utils.StatusMessage(fiber.StatusInternalServerError),
			"error":   err,
		})
	}
}
