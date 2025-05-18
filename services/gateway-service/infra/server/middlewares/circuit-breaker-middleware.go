package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kmlcnclk/kc-oms/common/pkg/handler"
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

		// Execute wrapped logic inside circuit breaker
		_, err := cb.Execute(func() (interface{}, error) {
			// Call next handler
			err := c.Next()

			// Log and consider server errors as breaker failures
			if c.Response().StatusCode() >= 500 {
				return nil, fiber.ErrInternalServerError
			}
			return nil, err
		})

		// Check if circuit is open or request failed
		if err != nil {
			if err == gobreaker.ErrOpenState {
				zap.L().Warn("Circuit breaker is open", zap.String("path", path))
				return c.Status(fiber.StatusServiceUnavailable).JSON(handler.ErrorResponse{
					Error: "Service temporarily unavailable. Please try again later.",
				},
				)
			}

			zap.L().Error("Circuit breaker execution failed", zap.String("path", path), zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(handler.ErrorResponse{
				Error: "Internal server error. Please try again later.",
			})
		}

		return nil
	}
}
