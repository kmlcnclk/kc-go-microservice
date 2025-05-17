package server

import (
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/kmlcnclk/kc-oms/services/gateway-service/infra/server/middlewares"
	"github.com/sony/gobreaker"
)

func Middleware(app *fiber.App, cb *gobreaker.CircuitBreaker) {
	app.Use(otelfiber.Middleware())

	app.Use(middlewares.CircuitBreakerMiddleware(cb, []string{}))

	app.Use(middlewares.LimiterMiddleware(100))
}
