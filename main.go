package main

import (
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Serve static files from the frontend build directory
	app.Static("/", filepath.Join("dist"))

	// Start the server on port 3000
	app.Listen(":3000")
}
