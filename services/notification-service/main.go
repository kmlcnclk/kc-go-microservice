package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	orderPB "github.com/kmlcnclk/kc-oms/common/api/order"
	"github.com/kmlcnclk/kc-oms/common/pkg/log"
	tracer "github.com/kmlcnclk/kc-oms/common/pkg/tracer"

	"github.com/kmlcnclk/kc-oms/common/pkg/config"
	"github.com/kmlcnclk/kc-oms/common/pkg/rabbitmq"
	notificationConfig "github.com/kmlcnclk/kc-oms/services/notification-service/infra/config"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var (
	serviceName = "notification-service"
	jaegerAddr  = "localhost:4318"
)

func main() {
	appConfig := config.ReadConfig[notificationConfig.AppConfig]()
	log.Init("[NOTIFICATION-SERVICE]")
	defer zap.L().Sync()

	_, err := tracer.SetGlobalTracer(context.Background(), serviceName, jaegerAddr)
	if err != nil {
		zap.L().Fatal("failed to set global tracer", zap.Error(err))
	}

	zap.L().Info("notification-service starting...")

	mq, err := rabbitmq.NewRabbitMQ(appConfig.RabbitMQURL)
	if err != nil {
		zap.L().Fatal("Could not connect to RabbitMQ", zap.Error(err))
	}
	defer mq.Close()

	if err := mq.Build(appConfig.RabbitOrderQueue, appConfig.RabbitOrderExchange, appConfig.RabbitOrderRoutingKey); err != nil {
		zap.L().Error("Failed to build queue/exchange", zap.Error(err))
	}

	zap.L().Info("RabbitMQ connected")

	messages, err := mq.Consume(appConfig.RabbitOrderQueue, false)
	if err != nil {
		zap.L().Fatal("Failed to start consuming messages", zap.Error(err))
	}

	tp := otel.Tracer("orders")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		for msg := range messages {
			zap.L().Info("Received message", zap.String("body", string(msg.Body)))

			messageCtx := tracer.ExtractTraceFromAMQPHeaders(msg.Headers)
			_, span := tp.Start(messageCtx, "NotificationService.CreateOrder")

			var order orderPB.CreateOrderResponse
			if err := json.Unmarshal(msg.Body, &order); err != nil {
				zap.L().Error("Failed to parse message", zap.Error(err))
				msg.Nack(false, false)
				continue
			}

			zap.L().Info("Processing order",
				zap.String("orderId", order.OrderId),
				zap.String("status", order.Status),
			)

			success := true // TODO: implement actual logic

			if success {
				msg.Ack(false)
			} else {
				zap.L().Warn("Order processing failed, sending to DLQ")
				msg.Nack(false, false)
			}

			span.End()
		}
	}()

	<-ctx.Done()
	zap.L().Info("Received shutdown signal, exiting...")
}
