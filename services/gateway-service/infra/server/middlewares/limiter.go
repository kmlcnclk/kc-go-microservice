package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/kmlcnclk/kc-oms/common/pkg/handler"
)

func LimiterMiddleware(maxRequests int) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        maxRequests,      // max 5 requests
		Expiration: 10 * time.Second, // per 10 seconds
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(handler.ErrorResponse{
				Error: "Too many requests",
			})
		},
	})
}
