package handlers

import (
	"fiber-gorm/internal/models"
	"fiber-gorm/internal/services"
	"fiber-gorm/internal/validators"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	Svc *services.UserService
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var payload models.CreateUserPayload

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Validate the payload using custom validator
	if err := validators.ValidateUserCreation(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": validators.FormatValidationError(err),
		})
	}

	user := models.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
		Hobby:    payload.Hobby,
	}

	if err := h.Svc.CreateUser(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}
