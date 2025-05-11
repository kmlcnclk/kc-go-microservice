package service

import (
	"context"

	omspb "github.com/kmlcnclk/kc-oms/common/api"
	rabbitmq "github.com/kmlcnclk/kc-oms/common/pkg/rabbitmq"
	"go.uber.org/zap"
)

type OrderService struct {
	rabbitmq       *rabbitmq.RabbitMQ
	mqExchangeName string
	mqRoutingKey   string
	// Add any dependencies or configurations needed for the service
}

func NewOrderService(rabbitmq *rabbitmq.RabbitMQ, mqExchangeName, mqRoutingKey string) *OrderService {
	return &OrderService{
		rabbitmq:       rabbitmq,
		mqExchangeName: mqExchangeName,
		mqRoutingKey:   mqRoutingKey,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *omspb.CreateOrderRequest) (*omspb.CreateOrderResponse, error) {
	// TODO: Implement the order creation logic here

	err := s.rabbitmq.Publish(
		s.mqExchangeName,
		s.mqRoutingKey,
		[]byte(`{"OrderId": "123456", "Status": "created"}`),
		"application/json",
		ctx,
	)
	if err != nil {
		zap.L().Info("Failed to publish message to RabbitMQ", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Order created and message published to RabbitMQ")

	return &omspb.CreateOrderResponse{
		OrderId: "123456",
		Status:  "created",
	}, nil
}
