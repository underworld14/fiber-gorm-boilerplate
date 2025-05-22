package middleware

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type contextKey string

const requestIDKey contextKey = "requestID"

func RequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := uuid.New().String()

		c.Locals(requestIDKey, requestID)
		c.Set("X-Request-ID", requestID)

		ctx := context.WithValue(c.Context(), requestIDKey, requestID)
		c.SetUserContext(ctx)

		return c.Next()
	}
}

func GetRequestID(c *fiber.Ctx) string {
	if id, ok := c.Locals(requestIDKey).(string); ok {
		return id
	}

	return ""
}

func GetRequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}

	return ""
}
