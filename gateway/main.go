package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kmlcnclk/kc-oms/common/pkg/config"
	_ "github.com/kmlcnclk/kc-oms/common/pkg/log"
	"github.com/kmlcnclk/kc-oms/gateway/app/healthcheck"
	"github.com/kmlcnclk/kc-oms/gateway/app/order"
	"github.com/kmlcnclk/kc-oms/gateway/infra/server"
	"go.uber.org/zap"
)

func main() {
	appConfig := config.Read()
	defer zap.L().Sync()

	zap.L().Info("app starting...")

	app := fiber.New(fiber.Config{
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		Concurrency:  256 * 1024,
	})

	healthcheckHandler := healthcheck.NewHealthCheckHandler()
	orderCreateHandler := order.NewCreateOrderHandler()

	server.InitRouters(app, healthcheckHandler, orderCreateHandler)

	server.Start(app, appConfig)

	server.GracefulShutdown(app)
}
