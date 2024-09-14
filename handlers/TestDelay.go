package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func TestDelay(c *fiber.Ctx) error {
	return c.JSON(map[string]interface{}{
		"message": "delayed",
	})
}
