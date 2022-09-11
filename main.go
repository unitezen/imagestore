package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/unitezen/imagestore/core"
	"github.com/unitezen/imagestore/handlers"
)

func main() {

	// Initialize app with custom error handler
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {

			// Defaults to status code 500 for non-http errors
			code := fiber.StatusInternalServerError

			// Attempt to use HTTP error code when applicable
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			// Set content-type to application/jSON
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

			// Return JSON error message
			return c.Status(code).JSON(&fiber.Map{"error": err.Error()})
		},
	})

	// Connect to the database
	core.ConnectDatabase()

	// Public API for registration and login
	app.Post("/users", handlers.AddUser)
	app.Post("/sessions", handlers.LoginUser)

	// Private API
	authenticatedEndpoints := app.Group("/", handlers.APIKeyMiddleware)
	authenticatedEndpoints.Delete("/sessions", handlers.LogoutUser)

	// Image upload/listing
	authenticatedEndpoints.Post("/images", handlers.UploadImage)
	authenticatedEndpoints.Get("/images", handlers.GetImages)
	authenticatedEndpoints.Get("/images/:id", handlers.GetImage)
	authenticatedEndpoints.Patch("/images/:id", handlers.PatchImageName)
	authenticatedEndpoints.Delete("/images/:id", handlers.DeleteImage)

	// Start the server
	app.Listen("0.0.0.0:8080")
}
