package redis

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type CacheConfig struct {
	Key        string
	NewMessage func() proto.Message
}

func RedisCacheInterceptor(rdb *Redis, ttl time.Duration, cacheable map[string]CacheConfig) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		zap.L().Info("Redis cache interceptor: received request", zap.String("method", info.FullMethod))

		conf, ok := cacheable[info.FullMethod]
		if !ok {
			return handler(ctx, req)
		}

		zap.L().Info("Redis cache interceptor: checking cache", zap.String("key", conf.Key))

		cachedVal, err := rdb.Get(ctx, conf.Key)
		if err == nil {
			msg := conf.NewMessage()
			if err := proto.Unmarshal([]byte(cachedVal), msg); err == nil {
				zap.L().Info("Cache hit", zap.String("key", conf.Key))
				return msg, nil
			}
			zap.L().Warn("Failed to unmarshal cached protobuf", zap.Error(err))
		}

		zap.L().Info("Cache miss", zap.String("key", conf.Key))
		resp, err := handler(ctx, req)
		if err != nil {
			return resp, err
		}

		zap.L().Info("Handler succeeded, caching response", zap.String("key", conf.Key))

		if pbResp, ok := resp.(proto.Message); ok {
			if data, err := proto.Marshal(pbResp); err == nil {
				if err := rdb.Set(ctx, conf.Key, data, ttl); err != nil {
					zap.L().Warn("Failed to set cache", zap.Error(err))

				} else {
					zap.L().Info("Successfully cached response", zap.String("key", conf.Key))
				}
			} else {
				zap.L().Warn("Failed to marshal response", zap.Error(err))
			}
		}

		return resp, nil
	}
}
