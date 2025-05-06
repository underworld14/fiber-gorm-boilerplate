package handlers

import (
	"fiber-gorm/internal/models"
	"fiber-gorm/internal/services"
	"fiber-gorm/internal/validators"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// TokenResponse represents the response for token generation endpoints
type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// AuthHandler handles authentication routes
type AuthHandler struct {
	AuthSvc *services.AuthService
}

func NewAuthHandler(authSvc *services.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthSvc: authSvc,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	// Parse the request body
	var payload models.CreateUserPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Validate the payload
	if errors := validators.ValidateUserCreation(&payload); errors != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": validators.FormatValidationError(errors, payload),
		})
	}

	// Register the user
	user, err := h.AuthSvc.RegisterUser(&payload)
	if err != nil {
		log.Error().Err(err).Msg("Failed to register user")

		// Check for specific errors to return appropriate status codes
		if err.Error() == "email already in use" {
			return c.Status(http.StatusConflict).JSON(fiber.Map{
				"error": "Email is already registered",
			})
		}

		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Generate tokens for the newly registered user
	accessToken, refreshToken, err := h.AuthSvc.CreateTokens(user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate authentication tokens",
		})
	}

	// Calculate token expiration (15 minutes from now)
	expiresAt := time.Now().Add(15 * time.Minute)

	// Return the tokens
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"token": TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "bearer",
			ExpiresAt:    expiresAt,
		},
		"user": user,
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// Parse the request body
	var payload models.LoginUserPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Authenticate the user
	user, accessToken, refreshToken, err := h.AuthSvc.LoginUser(&payload)
	if err != nil {
		log.Debug().Err(err).Str("email", payload.Email).Msg("Login failed")

		// For security reasons, don't specify whether email or password is incorrect
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Calculate token expiration (15 minutes from now)
	expiresAt := time.Now().Add(15 * time.Minute)

	// Return the tokens
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"token": TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "bearer",
			ExpiresAt:    expiresAt,
		},
		"user": user,
	})
}

// RefreshToken handles token refresh requests
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	// Parse the request body
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Validate the refresh token and generate new tokens
	accessToken, refreshToken, err := h.AuthSvc.RefreshTokens(req.RefreshToken)
	if err != nil {
		log.Debug().Err(err).Msg("Token refresh failed")
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired refresh token",
		})
	}

	// Calculate token expiration (15 minutes from now)
	expiresAt := time.Now().Add(15 * time.Minute)

	// Return the new tokens
	return c.Status(http.StatusOK).JSON(TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "bearer",
		ExpiresAt:    expiresAt,
	})
}

// Me returns the authenticated user's information
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	// Get the user ID from the context (set by auth middleware)
	userID := c.Locals("userID").(string)
	if userID == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not authenticated",
		})
	}

	// Retrieve the user from the database
	user, err := h.AuthSvc.UserRepo.FindUserById(userID)
	if err != nil {
		log.Error().Err(err).Str("userID", userID).Msg("Failed to retrieve user")
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Return user info (excluding sensitive data)
	return c.Status(http.StatusOK).JSON(user)
}

// Logout doesn't actually invalidate tokens (since they're stateless)
// In a production app, you'd implement token blacklisting or use Redis for token management
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}
