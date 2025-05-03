package main

import (
	"net"

	"github.com/kmlcnclk/kc-oms/common/pkg/config"
	_ "github.com/kmlcnclk/kc-oms/common/pkg/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/kmlcnclk/kc-oms/orders/app"
	orderConfig "github.com/kmlcnclk/kc-oms/orders/infra/config"
	"github.com/kmlcnclk/kc-oms/orders/service"
)

func main() {
	appConfig := config.ReadConfig[orderConfig.AppConfig]()
	defer zap.L().Sync()

	zap.L().Info("app starting...")

	grpcServer := grpc.NewServer()

	listener, err := net.Listen("tcp", ":"+appConfig.Port)
	if err != nil {
		zap.L().Fatal("Failed to create TCP listener: ", zap.Error(err))
	}

	zap.L().Info("gRPC server listening on port: ", zap.String("port", appConfig.Port))

	service := service.NewOrderService()

	app.NewGrpcHandler(grpcServer, service)

	defer listener.Close()

	if err := grpcServer.Serve(listener); err != nil {
		zap.L().Fatal("Failed to serve gRPC server: ", zap.Error(err))
	}

}
