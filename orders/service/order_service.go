package service

import (
	"context"

	omspb "github.com/kmlcnclk/kc-oms/common/api"
)

type OrderService struct {
	// Add any dependencies or configurations needed for the service
}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *omspb.CreateOrderRequest) (*omspb.CreateOrderResponse, error) {
	// TODO: Implement the order creation logic here

	return &omspb.CreateOrderResponse{
		OrderId: "123456",
		Status:  "asd",
	}, nil
}
