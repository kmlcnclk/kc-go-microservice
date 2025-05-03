package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	omspb "github.com/kmlcnclk/kc-oms/common/api"
	"github.com/kmlcnclk/kc-oms/common/pkg/config"
	_ "github.com/kmlcnclk/kc-oms/common/pkg/log"
	"github.com/kmlcnclk/kc-oms/gateway/app/healthcheck"
	"github.com/kmlcnclk/kc-oms/gateway/app/order"
	gatewayConfig "github.com/kmlcnclk/kc-oms/gateway/infra/config"
	"github.com/kmlcnclk/kc-oms/gateway/infra/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	appConfig := config.ReadConfig[gatewayConfig.AppConfig]()
	defer zap.L().Sync()

	zap.L().Info("app starting...")

	app := fiber.New(fiber.Config{
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		Concurrency:  256 * 1024,
	})

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(appConfig.OrderServiceAddress, opts...)
	if err != nil {
		zap.L().Fatal("failed to connect to order service", zap.Error(err))
	}
	defer conn.Close()

	orderServiceClient := omspb.NewOrderServiceClient(conn)
	orderCreateHandler := order.NewCreateOrderHandler(orderServiceClient)

	healthcheckHandler := healthcheck.NewHealthCheckHandler()

	server.InitRouters(app, healthcheckHandler, orderCreateHandler)

	server.Start(app, appConfig)

	server.GracefulShutdown(app)
}
