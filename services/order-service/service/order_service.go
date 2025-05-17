package service

import (
	"context"

	pb "github.com/kmlcnclk/kc-oms/common/api/order"
	rabbitmq "github.com/kmlcnclk/kc-oms/common/pkg/rabbitmq"
	"go.opentelemetry.io/otel"
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

func (s *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	// TODO: Implement the order creation logic here

	tr := otel.Tracer("orders")
	ctx, messageSpan := tr.Start(ctx, "OrdersService.CreateOrder")
	defer messageSpan.End()

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

	return &pb.CreateOrderResponse{
		OrderId: "123456",
		Status:  "created",
	}, nil
}
