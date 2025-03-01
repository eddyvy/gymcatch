package main

import (
	"path/filepath"

	backend "github.com/eddyvy/gymcatch/backend"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	backend.InitSessions()

	app := fiber.New()

	// Handle login
	app.Post("/api/auth", backend.HandleAuth)
	app.Get("/api/check_session/:session", backend.HandleCheckSession)

	api := app.Group("/api", backend.SessionMiddleware)
	api.Get("/mega_events", backend.HandleMegaEvents)
	api.Post("/mega_inscribe/:classId", backend.HandleInscribe)
	api.Get("/mega_inscribe", backend.HandleGetInscribedClasses)

	// Serve static files from the frontend build directory
	app.Static("/", filepath.Join("dist"))

	// Start the server on port 3000
	app.Listen(":3000")
}
