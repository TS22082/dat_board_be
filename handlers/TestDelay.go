package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// TestDelay is a handler function that returns a JSON response with a message.
// The message is "delayed".
//
// Parameters: None
//
// Returns: JSON response with a message
func TestDelay(c *fiber.Ctx) error {
	return c.JSON(map[string]interface{}{
		"message": "delayed",
	})
}
