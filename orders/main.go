package main

import (
	"context"
	"net"
	"time"

	"github.com/kmlcnclk/kc-oms/common/pkg/config"
	_ "github.com/kmlcnclk/kc-oms/common/pkg/log"
	"github.com/kmlcnclk/kc-oms/common/pkg/rabbitmq"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/kmlcnclk/kc-oms/common/pkg/discovery"
	"github.com/kmlcnclk/kc-oms/common/pkg/discovery/consul"
	"github.com/kmlcnclk/kc-oms/orders/app"
	orderConfig "github.com/kmlcnclk/kc-oms/orders/infra/config"
	"github.com/kmlcnclk/kc-oms/orders/service"
)

var (
	serviceName = "orders"
	consulAddr  = "localhost:8500"
	grpcAddr    = "localhost:50052"
)

func main() {
	appConfig := config.ReadConfig[orderConfig.AppConfig]()
	defer zap.L().Sync()

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

	grpcServer := grpc.NewServer()

	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		zap.L().Fatal("Failed to create TCP listener: ", zap.Error(err))
	}

	rmq, err := rabbitmq.NewRabbitMQ(appConfig.RabbitMQURL)
	if err != nil {
		zap.L().Fatal("Failed to create RabbitMQ instance: ", zap.Error(err))
	}

	if err := rmq.Build(appConfig.RabbitMQQueueName, appConfig.RabbitMQExchangeName, appConfig.RabbitMQRoutingKey); err != nil {
		zap.L().Error("Failed to build queue/exchange", zap.Error(err))
	}

	zap.L().Info("gRPC server listening on port: ", zap.String("port", grpcAddr))

	service := service.NewOrderService(rmq, appConfig.RabbitMQExchangeName, appConfig.RabbitMQRoutingKey)

	app.NewGrpcHandler(grpcServer, service)

	defer listener.Close()

	if err := grpcServer.Serve(listener); err != nil {
		zap.L().Fatal("Failed to serve gRPC server: ", zap.Error(err))
	}

}
