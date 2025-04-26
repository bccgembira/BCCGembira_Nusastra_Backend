package response

import (
	"github.com/gofiber/fiber/v2"
)

func Success(c *fiber.Ctx, key string, data interface{}) error {
	response := map[string]interface{}{
		"message": "success",
	}

	if key != "" {
		response[key] = data
	} else {
		response["detail"] = data
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
