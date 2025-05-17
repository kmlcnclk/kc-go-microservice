package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"go.uber.org/zap"
)

type Redis struct {
	client *redis.Client
	tracer trace.Tracer
}

func NewRedis(addr, password string) *Redis {
	ctx := context.Background()
	zap.L().Info("Initializing Redis client", zap.String("addr", addr))
	tracer := otel.Tracer("redis")

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		zap.L().Error("Failed to connect to Redis", zap.Error(err))
	}

	zap.L().Info("Connected to Redis successfully")

	return &Redis{
		client: rdb,
		tracer: tracer,
	}
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	ctx, span := r.tracer.Start(ctx, "Redis SET")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.redis.key", key),
		attribute.String("db.redis.ttl", expiration.String()),
	)

	err := r.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		zap.L().Error("Failed to set value in Redis", zap.String("key", key), zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	ctx, span := r.tracer.Start(ctx, "Redis GET")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.redis.key", key),
	)

	return r.client.Get(ctx, key).Result()
}

func (r *Redis) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
