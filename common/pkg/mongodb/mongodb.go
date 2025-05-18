package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
	Ctx      context.Context
	Cancel   context.CancelFunc
}

func NewMongoDB(uri, dbName string) *MongoDB {
	tracer := otel.Tracer("mongo")

	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, evt *event.CommandStartedEvent) {
			_, span := tracer.Start(ctx, evt.CommandName)
			zap.L().Info("MongoDB command started", zap.String("command", evt.CommandName))
			span.End()
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().
		ApplyURI(uri).
		SetMonitor(monitor)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		cancel()
		zap.L().Fatal("MongoDB connection error", zap.Error(err))
	}

	if err = client.Ping(ctx, nil); err != nil {
		cancel()
		zap.L().Fatal("MongoDB ping error", zap.Error(err))
	}

	zap.L().Info("MongoDB connection established", zap.String("uri", uri))

	return &MongoDB{
		Client:   client,
		Database: client.Database(dbName),
		Ctx:      ctx,
		Cancel:   cancel,
	}
}

func (m *MongoDB) GetCollection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}

func (m *MongoDB) Close() {
	m.Cancel()
	if err := m.Client.Disconnect(m.Ctx); err != nil {
		zap.L().Error("MongoDB disconnect error", zap.Error(err))
	} else {
		zap.L().Info("MongoDB disconnected")
	}
}
