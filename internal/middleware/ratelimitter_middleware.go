package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type RateLimitter struct {
	mu          sync.Mutex
	requests    map[string]int
	limit       int
	resetPeriod time.Duration
}

func NewRateLimitter(limit int, resetPeriod time.Duration) *RateLimitter {
	return &RateLimitter{
		requests:    make(map[string]int),
		limit:       limit,
		resetPeriod: resetPeriod,
	}
}

func (rl *RateLimitter) Middleware(c *fiber.Ctx) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	ip := c.IP()
	count := rl.requests[ip]

	if count >= rl.limit {
		return fiber.NewError(fiber.StatusTooManyRequests, "Too Many Request")
	}

	rl.requests[ip]++

	time.AfterFunc(rl.resetPeriod, func() {
		rl.mu.Lock()
		defer rl.mu.Unlock()
		delete(rl.requests, ip)
	})

	return c.Next()
}
