package server

import (
	"github.com/gofiber/fiber/v2"
	orderPB "github.com/kmlcnclk/kc-oms/common/api/order"
	productPB "github.com/kmlcnclk/kc-oms/common/api/product"
	"github.com/kmlcnclk/kc-oms/common/pkg/handler"
	"github.com/kmlcnclk/kc-oms/services/gateway-service/app/healthcheck"
	"github.com/kmlcnclk/kc-oms/services/gateway-service/app/order"
	"github.com/kmlcnclk/kc-oms/services/gateway-service/app/product"
)

func InitRouters(app *fiber.App, healthcheckHandler *healthcheck.HealthCheckHandler, orderCreateHandler *order.CreateOrderHandler, productGetAllProductsHandler *product.GetAllProductsHandler) {

	app.Get("/healthcheck", handler.Handle[healthcheck.HealthCheckRequest, healthcheck.HealthCheckResponse](healthcheckHandler))

	app.Post("/order", handler.Handle[orderPB.CreateOrderRequest, handler.SuccessResponse](orderCreateHandler))

	app.Get("/product", handler.Handle[productPB.GetAllProductsRequest, handler.SuccessResponse](productGetAllProductsHandler))

}
