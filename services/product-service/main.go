package main

import (
	"context"
	"net"
	"time"

	pb "github.com/kmlcnclk/kc-oms/common/api/product"
	"github.com/kmlcnclk/kc-oms/common/pkg/config"
	"github.com/kmlcnclk/kc-oms/common/pkg/discovery"
	"github.com/kmlcnclk/kc-oms/common/pkg/discovery/consul"
	"github.com/kmlcnclk/kc-oms/common/pkg/log"
	"github.com/kmlcnclk/kc-oms/common/pkg/mongodb"
	"github.com/kmlcnclk/kc-oms/common/pkg/redis"
	tracer "github.com/kmlcnclk/kc-oms/common/pkg/tracer"
	"github.com/kmlcnclk/kc-oms/services/product-service/app"
	productConfig "github.com/kmlcnclk/kc-oms/services/product-service/infra/config"
	"github.com/kmlcnclk/kc-oms/services/product-service/repositories"
	"github.com/kmlcnclk/kc-oms/services/product-service/services"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

var (
	serviceName = "products"
	consulAddr  = "localhost:8500"
	grpcAddr    = "localhost:50052"
	jaegerAddr  = "localhost:4318"
)

func main() {
	appConfig := config.ReadConfig[productConfig.AppConfig]()
	log.Init("[PRODUCT-SERVICE]")
	defer zap.L().Sync()

	tp, err := tracer.SetGlobalTracer(context.TODO(), serviceName, jaegerAddr)
	if err != nil {
		zap.L().Fatal("failed to set global tracer", zap.Error(err))
	}

	zap.L().Info("app starting...")

	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		panic(err)
	}

	defer registry.Deregister(ctx, instanceID, serviceName)

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				zap.L().Error("Failed to health check", zap.Error(err))
			}
			time.Sleep(time.Second * 1)
		}
	}()

	mongoDB := mongodb.NewMongoDB(appConfig.MONGO_URI, appConfig.MONGO_DB)

	rdb := redis.NewRedis(appConfig.REDIS_ADDR, appConfig.REDIS_PASS)

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(
			otelgrpc.NewServerHandler(
				otelgrpc.WithTracerProvider(tp),
			),
		),
		grpc.UnaryInterceptor(
			redis.RedisCacheInterceptor(rdb, time.Hour, map[string]redis.CacheConfig{
				"/api.ProductService/GetAllProducts": {
					Key: "all_products",
					NewMessage: func() proto.Message {
						return &pb.GetAllProductsResponse{}
					},
				},
			}),
		),
	)

	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		zap.L().Fatal("Failed to create TCP listener: ", zap.Error(err))
	}

	zap.L().Info("gRPC server listening on port: ", zap.String("port", grpcAddr))

	repository := repositories.NewProductRepository(mongoDB, appConfig.MONGO_COLLECTION)
	service := services.NewProductService(repository)

	app.NewGrpcHandler(grpcServer, service)

	defer listener.Close()

	if err := grpcServer.Serve(listener); err != nil {
		zap.L().Fatal("Failed to serve gRPC server: ", zap.Error(err))
	}

}
