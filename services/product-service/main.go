package main

import (
	"context"
	"net"
	"time"

	"github.com/kmlcnclk/kc-oms/common/pkg/config"
	"github.com/kmlcnclk/kc-oms/common/pkg/log"
	tracer "github.com/kmlcnclk/kc-oms/common/pkg/tracer"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/kmlcnclk/kc-oms/common/pkg/discovery"
	"github.com/kmlcnclk/kc-oms/common/pkg/discovery/consul"
	"github.com/kmlcnclk/kc-oms/services/product-service/app"
	productConfig "github.com/kmlcnclk/kc-oms/services/product-service/infra/config"
	"github.com/kmlcnclk/kc-oms/services/product-service/service"
)

var (
	serviceName = "products"
	consulAddr  = "localhost:8500"
	grpcAddr    = "localhost:50052"
	jaegerAddr  = "localhost:4318"
)

func main() {
	_ = config.ReadConfig[productConfig.AppConfig]()
	log.Init("[PRODUCT-SERVICE]")
	defer zap.L().Sync()

	tp, err := tracer.SetGlobalTracer(context.TODO(), serviceName, jaegerAddr)
	if err != nil {
		zap.L().Fatal("failed to set global tracer", zap.Error(err))
	}

	zap.L().Info("app starting...")

	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		panic(err)
	}

	defer registry.Deregister(ctx, instanceID, serviceName)

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				zap.L().Error("Failed to health check", zap.Error(err))
			}
			time.Sleep(time.Second * 1)
		}
	}()

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(
			otelgrpc.NewServerHandler(
				otelgrpc.WithTracerProvider(tp),
			),
		),
	)

	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		zap.L().Fatal("Failed to create TCP listener: ", zap.Error(err))
	}

	zap.L().Info("gRPC server listening on port: ", zap.String("port", grpcAddr))

	service := service.NewProductService()

	app.NewGrpcHandler(grpcServer, service)

	defer listener.Close()

	if err := grpcServer.Serve(listener); err != nil {
		zap.L().Fatal("Failed to serve gRPC server: ", zap.Error(err))
	}

}
