package main

import (
	"fmt"
	// "path/filepath"

	backend "github.com/eddyvy/gymcatch/backend"
	// "github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	// sessions := backend.NewSessions()

	sId, csrfToken, err := backend.GetCreds()
	if err != nil {
		panic(err)
	}

	fmt.Println("Session ID:", sId)
	fmt.Println("CSRF Token:", csrfToken)

	// app := fiber.New()

	// // Handle login
	// app.Post("/api/auth", backend.HandleAuth(sessions))

	// // Serve static files from the frontend build directory
	// app.Static("/", filepath.Join("dist"))

	// // Start the server on port 3000
	// app.Listen(":3000")
}
