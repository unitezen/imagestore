package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/unitezen/imagestore/core"
	"github.com/unitezen/imagestore/models"
)

func APIKeyMiddleware(c *fiber.Ctx) error {

	// Extract api key from headers
	api_key := c.Get("X-API-Key")
	if api_key == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "X-API-Key header required for accessing private endpoints")
	}

	// Retrieve session
	userSession := &models.UserSession{ApiKey: api_key}
	result := core.Database.First(&userSession)

	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "API key not found")
	}

	user := &models.User{ID: userSession.UserId}
	userResult := core.Database.First(&user)
	if userResult.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	c.Locals("user", user)
	return c.Next()
}
