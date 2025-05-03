package server

import (
	"github.com/gofiber/fiber/v2"
	omspb "github.com/kmlcnclk/kc-oms/common/api"
	"github.com/kmlcnclk/kc-oms/common/pkg/handler"
	"github.com/kmlcnclk/kc-oms/gateway/app/healthcheck"
	"github.com/kmlcnclk/kc-oms/gateway/app/order"
)

func InitRouters(app *fiber.App, healthcheckHandler *healthcheck.HealthCheckHandler, orderCreateHandler *order.CreateOrderHandler) {

	app.Get("/healthcheck", handler.Handle[healthcheck.HealthCheckRequest, healthcheck.HealthCheckResponse](healthcheckHandler))

	app.Post("/order", handler.Handle[omspb.CreateOrderRequest, omspb.CreateOrderResponse](orderCreateHandler))

}
