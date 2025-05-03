package app

import (
	"context"

	omspb "github.com/kmlcnclk/kc-oms/common/api"
	"github.com/kmlcnclk/kc-oms/orders/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GrpcHandler struct {
	omspb.UnimplementedOrderServiceServer
	service *service.OrderService
}

func NewGrpcHandler(grpcServer *grpc.Server, service *service.OrderService) {
	handler := &GrpcHandler{service: service}
	zap.L().Info("Registering CreateOrderHandler with gRPC server")
	omspb.RegisterOrderServiceServer(grpcServer, handler)
	zap.L().Info("CreateOrderHandler registered successfully")

}

func (g *GrpcHandler) CreateOrder(ctx context.Context, req *omspb.CreateOrderRequest) (*omspb.CreateOrderResponse, error) {
	zap.L().Info("Creating order")

	// TODO: Process the order creation logic here
	// For example, you can call a service to create the order

	response, err := g.service.CreateOrder(ctx, req)
	if err != nil {
		zap.L().Error("Failed to create order", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Order created successfully", zap.String("order_id", response.OrderId))

	return response, nil
}
