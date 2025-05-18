package main

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	circuitbreaker "github.com/kmlcnclk/kc-oms/common/pkg/circuit-breaker"
	"github.com/kmlcnclk/kc-oms/common/pkg/config"
	"github.com/kmlcnclk/kc-oms/common/pkg/log"
	tracer "github.com/kmlcnclk/kc-oms/common/pkg/tracer"
	"github.com/kmlcnclk/kc-oms/services/gateway-service/app/healthcheck"
	"github.com/kmlcnclk/kc-oms/services/gateway-service/app/order"
	"github.com/kmlcnclk/kc-oms/services/gateway-service/app/product"
	gatewayConfig "github.com/kmlcnclk/kc-oms/services/gateway-service/infra/config"
	"github.com/kmlcnclk/kc-oms/services/gateway-service/infra/server"
	"go.uber.org/zap"

	"github.com/kmlcnclk/kc-oms/common/pkg/discovery"
	"github.com/kmlcnclk/kc-oms/common/pkg/discovery/consul"
)

var (
	serviceName = "gateway"
	httpAddr    = ":8080"
	consulAddr  = "localhost:8500"
	jaegerAddr  = "localhost:4318"
)

func main() {
	appConfig := config.ReadConfig[gatewayConfig.AppConfig]()
	log.Init("[GATEWAY-SERVICE]")
	defer zap.L().Sync()

	tp, err := tracer.SetGlobalTracer(context.TODO(), serviceName, jaegerAddr)
	if err != nil {
		zap.L().Fatal("failed to set global tracer", zap.Error(err))
	}

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
				zap.L().Error("failed to health check", zap.Error(err))
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

	orderCreateHandler := order.NewCreateOrderHandler(registry, tp)
	productGetAllProductsHandler := product.NewGetAllProductsHandler(registry, tp)

	healthcheckHandler := healthcheck.NewHealthCheckHandler()

	cb := circuitbreaker.DefaultCircuitBreaker()

	server.Middleware(app, cb)

	server.InitRouters(app, healthcheckHandler, orderCreateHandler, productGetAllProductsHandler)

	server.Start(app, appConfig)

	server.GracefulShutdown(app)
}
