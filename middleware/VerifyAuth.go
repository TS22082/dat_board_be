package middleware

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func VerifyAuth(c *fiber.Ctx) error {
	authToken := c.Get("Authorization")

	if authToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "did not receive a valid token",
		})
	}

	// Parse the token
	parsedToken, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token's algorithm matches the expected algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "failed to parse token",
		})
	}

	// Validate token and extract claims
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		// Check the expiration time
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "token is expired",
				})
			}
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token claims",
			})
		}

		// Token is valid and not expired
		userID := claims["id"]
		c.Locals("userId", userID)

		return c.Next()
	}

	// Token is invalid
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "invalid token",
	})
}
