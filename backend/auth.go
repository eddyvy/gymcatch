package backend

import (
	"crypto/rand"
	"encoding/hex"
	"os"

	"github.com/gofiber/fiber/v2"
)

// UserCredentials represents the structure of the login request body
type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// generateSessionID generates a random session ID
func generateSessionID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HandleAuth handles the login request
func HandleAuth(c *fiber.Ctx) error {
	var creds UserCredentials

	// Parse the JSON body
	if err := c.BodyParser(&creds); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	// Check user credentials (this is just an example, you should replace it with your own logic)
	if creds.Email != "" && creds.Password == os.Getenv("PASSWORD") {
		// Generate a session ID
		sessionID, err := generateSessionID()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to generate session ID",
			})
		}

		// New credentials
		credsMega := NewMegaCreds(creds.Email)
		credsMega.LoadCreds()

		// Save the email into the session service
		Sessions.Set(sessionID, credsMega)

		// Return the session ID
		return c.JSON(fiber.Map{
			"sessionID": sessionID,
		})
	}

	// If the credentials are incorrect, return an unauthorized status
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "invalid credentials",
	})
}

func HandleCheckSession(c *fiber.Ctx) error {
	sessionID := c.Params("session")
	_, exists := Sessions.Get(sessionID)

	return c.JSON(fiber.Map{
		"success": exists,
	})
}
