package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// RequestLogger logs information about each request
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		path := c.Path()
		method := c.Method()

		// Process request
		err := c.Next()

		// Log request details after completion
		latency := time.Since(start)
		status := c.Response().StatusCode()

		log.Info().
			Str("method", method).
			Str("path", path).
			Int("status", status).
			Dur("latency", latency).
			Str("ip", c.IP()).
			// Str("user_agent", c.Get("User-Agent")).
			Msg("Request processed")

		return err
	}
}

// ErrorHandler provides a consistent error response format
func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// Default to 500 Internal Server Error
		code := fiber.StatusInternalServerError

		// Check if it's a Fiber error
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		// Log the error
		log.Error().
			Err(err).
			Int("status", code).
			Str("path", c.Path()).
			Str("method", c.Method()).
			Msg("Request error")

		// Return JSON error response
		return c.Status(code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
}
