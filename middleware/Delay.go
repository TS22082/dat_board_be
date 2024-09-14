package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Delay(c *fiber.Ctx) error {

	// get the delay parameter from the URL
	delay := c.Params("delay")

	// convert the delay parameter to an int
	delayInt, err := strconv.Atoi(delay)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid delay parameter",
		})
	}

	time.Sleep(time.Duration(delayInt) * time.Second)
	return c.Next()
}
