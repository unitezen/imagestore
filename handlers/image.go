package handlers

import (
	b64 "encoding/base64"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/unitezen/imagestore/core"
	"github.com/unitezen/imagestore/models"
)

func UploadImage(c *fiber.Ctx) error {

	// Get user from context
	user := c.Locals("user").(*models.User)

	// Marshal input payload into struct
	imageData := &models.ImageData{}
	if err := c.BodyParser(imageData); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Error parsing request payload")
	}

	// Validate struct
	err := models.ValidatePayload(imageData)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	// Verify if image is png/jpeg
	decodedImage, decodeErr := b64.StdEncoding.DecodeString(imageData.Data)
	if decodeErr != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Error decoding base64-encoded image data")
	}

	detectedType := http.DetectContentType(decodedImage)
	if (detectedType != "image/png") && (detectedType != "image/jpeg") {
		return fiber.NewError(fiber.StatusBadRequest, "Only PNG and JPEG images are allowed")
	}

	image := &models.Image{
		UserId:    user.ID,
		ImageData: *imageData,
	}

	result := core.Database.Create(&image)
	if result.RowsAffected == 1 {
		return c.Status(fiber.StatusOK).JSON(&fiber.Map{"id": image.ID})
	}
	return result.Error
}

func GetImage(c *fiber.Ctx) error {

	// Get user from context
	user := c.Locals("user").(*models.User)
	image := &models.Image{}

	// Fetch the user image
	result := core.Database.Where("id = ? AND user_id = ?", c.Params("id"), user.ID).First(&image)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Image with specified ID not found for this user")
	}
	return c.Status(fiber.StatusOK).JSON(image)

}

func DeleteImage(c *fiber.Ctx) error {

	// Get user from context
	user := c.Locals("user").(*models.User)
	image := &models.Image{}

	// Fetch the user image
	result := core.Database.Where("id = ? AND user_id = ?", c.Params("id"), user.ID).First(&image)
	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Image with specified ID not found for this user")
	}
	deleteResult := core.Database.Delete(&image)
	if deleteResult.Error != nil {
		return result.Error
	}
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"success": true})

}

func GetImages(c *fiber.Ctx) error {

	// Get user from context
	user := c.Locals("user").(*models.User)

	// Fetch the user image
	images := []models.Image{}
	result := core.Database.Where("user_id = ?", user.ID).Omit("data").Find(&images)
	if result.Error != nil {
		return result.Error
	}
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"images": &images})

}

func PatchImageName(c *fiber.Ctx) error {

	// Get user from context
	user := c.Locals("user").(*models.User)

	// Marshal input payload into struct
	imageName := &models.ImageName{}
	if err := c.BodyParser(imageName); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Error parsing request payload")
	}

	// Validate struct
	err := models.ValidatePayload(imageName)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	image := &models.Image{}
	// Fetch the user image
	result := core.Database.Where("id = ? AND user_id = ?", c.Params("id"), user.ID).First(&image)
	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Image with specified ID not found for this user")
	}

	// Update name
	image.Name = imageName.Name
	saveResult := core.Database.Save(&image)
	if saveResult.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Name update failed")
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"id": c.Params("id"), "name": imageName.Name})

}
