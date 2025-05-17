package service

import (
	"context"

	pb "github.com/kmlcnclk/kc-oms/common/api/product"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

type ProductService struct {
	// Add any dependencies or configurations needed for the service
}

func NewProductService() *ProductService {
	return &ProductService{}
}

func (s *ProductService) GetAllProducts(ctx context.Context, req *pb.GetAllProductsRequest) (*pb.GetAllProductsResponse, error) {
	// TODO: Implement the product fetching logic here

	tr := otel.Tracer("products")
	_, messageSpan := tr.Start(ctx, "ProductService.GetProducts")
	defer messageSpan.End()

	zap.L().Info("Products successfully returned!")

	return &pb.GetAllProductsResponse{
		Products: []*pb.Product{
			{
				Id:   "1",
				Name: "Product 1",
			},
		},
	}, nil
}
