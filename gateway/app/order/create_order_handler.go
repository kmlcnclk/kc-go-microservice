package order

import (
	"context"

	omspb "github.com/kmlcnclk/kc-oms/common/api"
	"go.uber.org/zap"
)

type CreateOrderHandler struct {
	client omspb.OrderServiceClient
}

func NewCreateOrderHandler(client omspb.OrderServiceClient) *CreateOrderHandler {
	return &CreateOrderHandler{client}
}

func (h *CreateOrderHandler) Handle(ctx context.Context, req *omspb.CreateOrderRequest) (*omspb.CreateOrderResponse, error) {
	zap.L().Info("Creating order", zap.String("order_id", req.CustomerId))

	createdOrder, err := h.client.CreateOrder(ctx, req)

	if err != nil {
		zap.L().Error("Failed to create order", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Order created successfully", zap.String("order_id", req.CustomerId))

	return createdOrder, nil
}
