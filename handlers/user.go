package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/unitezen/imagestore/core"
	"github.com/unitezen/imagestore/models"
	"github.com/unitezen/imagestore/utilities"
)

func AddUser(c *fiber.Ctx) error {

	// Marshal input payload into struct
	userData := &models.UserRegistration{}
	if err := c.BodyParser(userData); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Error parsing request payload")
	}

	// Validate struct
	err := models.ValidatePayload(userData)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	// Attempt to hash password
	hashedPassword, hashErr := utilities.HashPassword(userData.Password)
	if hashErr != nil {
		return hashErr
	}

	// Attempt to create user
	user := &models.User{
		Username: userData.Username,
		Password: hashedPassword,
		Email:    userData.Email,
	}
	result := core.Database.Create(user)

	if result.RowsAffected == 1 {
		// Create a new API key and return it
		userSession, err := utilities.CreateUserAPIKey(user)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(userSession)
	}

	return errors.New("Username/email already exists")
}

func LoginUser(c *fiber.Ctx) error {

	// Marshal input payload into struct
	userData := &models.UserLogin{}
	if err := c.BodyParser(userData); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Error parsing request payload")
	}

	// Validate struct
	err := models.ValidatePayload(userData)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	// Find the user
	user := &models.User{Username: userData.Username}
	result := core.Database.First(&user)
	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Username not found")
	}

	// Compare hash
	correctPassword := utilities.ComparePlaintextAndHash(user.Password, userData.Password)
	if !correctPassword {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid password")
	}

	// Delete all existing API key for user
	utilities.DeleteUserAPIKey(user)

	// Create new API key
	session, sessionErr := utilities.CreateUserAPIKey(user)
	if sessionErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, sessionErr.Error())
	}
	return c.Status(fiber.StatusOK).JSON(session)

}

func LogoutUser(c *fiber.Ctx) error {

	// Get authenticated user from context
	user := c.Locals("user").(*models.User)

	// Delete all existing API key for user
	utilities.DeleteUserAPIKey(user)

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"success": true})
}
