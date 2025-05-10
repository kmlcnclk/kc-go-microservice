package order

import (
	"context"

	omspb "github.com/kmlcnclk/kc-oms/common/api"
	"github.com/kmlcnclk/kc-oms/common/pkg/discovery"
	"go.uber.org/zap"
)

type CreateOrderHandler struct {
	registry discovery.Registry
}

func NewCreateOrderHandler(registry discovery.Registry) *CreateOrderHandler {
	return &CreateOrderHandler{
		registry: registry,
	}

}

func (h *CreateOrderHandler) Handle(ctx context.Context, req *omspb.CreateOrderRequest) (*omspb.CreateOrderResponse, error) {
	zap.L().Info("Creating order", zap.String("order_id", req.CustomerId))

	conn, err := discovery.ServiceConnection(context.Background(), "orders", h.registry)
	if err != nil {
		zap.L().Error("Failed to dial server", zap.Error(err))
	}
	defer conn.Close()

	createdOrder, err := omspb.NewOrderServiceClient(conn).CreateOrder(ctx, req)

	if err != nil {
		zap.L().Error("Failed to create order", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Order created successfully", zap.String("order_id", req.CustomerId))

	return createdOrder, nil
}
