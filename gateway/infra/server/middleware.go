package server

import (
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
)

func Middleware(app *fiber.App) {
	app.Use(otelfiber.Middleware())
}
