package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

func CircuitBreakerMiddleware(cb *gobreaker.CircuitBreaker, excludePaths []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()

		// Skip circuit breaker for excluded paths
		for _, excluded := range excludePaths {
			if path == excluded {
				return c.Next()
			}
		}
		zap.L().Info("Circuit breaker middleware", zap.String("path", path))
		// Wrap in circuit breaker
		// _, err := cb.Execute(func() (interface{}, error) {
		// 	return nil, c.Next()
		// })
		_, err := cb.Execute(func() (interface{}, error) {
			err := c.Next()

			// Consider specific status codes as failures
			if c.Response().StatusCode() >= 500 {
				return nil, fiber.ErrInternalServerError
			}

			return nil, err
		})

		return err
	}
}
