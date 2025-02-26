package backend

import (
	"github.com/gofiber/fiber/v2"
)

var sessionKey = "X-Session"

// SessionMiddleware checks if the X-Session header is present and validates the session
func SessionMiddleware(c *fiber.Ctx) error {
	sessionID := c.Get(sessionKey)
	if sessionID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing X-Session header",
		})
	}

	// Validate the session
	_, exists := Sessions.Get(sessionID)
	if !exists {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid session",
		})
	}

	// Proceed to the next handler
	return c.Next()
}

func GetSessionID(c *fiber.Ctx) string {
	return c.Get(sessionKey)
}
