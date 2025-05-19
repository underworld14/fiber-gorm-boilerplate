package middleware

import (
	"fiber-gorm/internal/config"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

// JWTAuthMiddleware creates a middleware for protecting routes with JWT
func JWTAuthMiddleware(cfg *config.Config) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.HS256,
			Key:    []byte(cfg.JWTSecret),
		},
		ContextKey:   "user",
		ErrorHandler: jwtError,
		TokenLookup:  "header:Authorization:Bearer ,query:token",
		SuccessHandler: func(c *fiber.Ctx) error {
			token := c.Locals("user").(*jwt.Token)
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				c.Locals("userID", claims["sub"]) // Set user ID in context from 'sub' claim
			}
			return c.Next()
		},
	})
}

// jwtError handles JWT validation errors
func jwtError(c *fiber.Ctx, err error) error {
	log.Error().Err(err).Msg("JWT validation error")

	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing or malformed token",
		})
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": "Invalid or expired token",
	})
}
