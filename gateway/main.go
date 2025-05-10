package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kmlcnclk/kc-oms/common/pkg/config"
	_ "github.com/kmlcnclk/kc-oms/common/pkg/log"
	"github.com/kmlcnclk/kc-oms/gateway/app/healthcheck"
	"github.com/kmlcnclk/kc-oms/gateway/app/order"
	gatewayConfig "github.com/kmlcnclk/kc-oms/gateway/infra/config"
	"github.com/kmlcnclk/kc-oms/gateway/infra/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/kmlcnclk/kc-oms/common/pkg/discovery"
	"github.com/kmlcnclk/kc-oms/common/pkg/discovery/consul"
)

var (
	serviceName = "gateway"
	httpAddr    = ":8080"
	consulAddr  = "localhost:50053"
)

func main() {
	appConfig := config.ReadConfig[gatewayConfig.AppConfig]()
	defer zap.L().Sync()

	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, httpAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatal("failed to health check")
			}
			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

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

	orderCreateHandler := order.NewCreateOrderHandler(registry)

	healthcheckHandler := healthcheck.NewHealthCheckHandler()

	server.InitRouters(app, healthcheckHandler, orderCreateHandler)

	server.Start(app, appConfig)

	server.GracefulShutdown(app)
}
