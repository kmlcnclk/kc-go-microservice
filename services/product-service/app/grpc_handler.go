package app

import (
	"context"

	pb "github.com/kmlcnclk/kc-oms/common/api/product"
	"github.com/kmlcnclk/kc-oms/services/product-service/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GrpcHandler struct {
	pb.UnimplementedProductServiceServer
	service *service.ProductService
}

func NewGrpcHandler(grpcServer *grpc.Server, service *service.ProductService) {
	handler := &GrpcHandler{service: service}
	zap.L().Info("Registering ProductHandler with gRPC server")
	pb.RegisterProductServiceServer(grpcServer, handler)
	zap.L().Info("ProductHandler registered successfully")

}

func (g *GrpcHandler) GetAllProducts(ctx context.Context, req *pb.GetAllProductsRequest) (*pb.GetAllProductsResponse, error) {
	zap.L().Info("All products requested")

	response, err := g.service.GetAllProducts(ctx, req)
	if err != nil {
		zap.L().Error("Failed to fetch all products", zap.Error(err))
		return nil, err
	}

	zap.L().Info("All Products successfully returned!")

	return response, nil
}
