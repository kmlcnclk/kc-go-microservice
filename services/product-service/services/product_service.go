package services

import (
	"context"

	pb "github.com/kmlcnclk/kc-oms/common/api/product"
	"github.com/kmlcnclk/kc-oms/common/pkg/models"
	"github.com/kmlcnclk/kc-oms/services/product-service/repositories"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type ProductServiceInterface interface {
	GetAllProducts(ctx context.Context) (*pb.GetAllProductsResponse, error)
}

type ProductService struct {
	repo   *repositories.ProductRepository
	tracer trace.Tracer
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	tracer := otel.Tracer("products")
	return &ProductService{
		repo:   repo,
		tracer: tracer,
	}
}

func (s *ProductService) GetAllProducts(ctx context.Context) (*pb.GetAllProductsResponse, error) {
	_, messageSpan := s.tracer.Start(ctx, "ProductService.GetProducts")
	defer messageSpan.End()

	res, err := s.repo.GetAllProducts(ctx)
	if err != nil {
		zap.L().Error("Failed to fetch all products", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Products successfully returned!")

	return toGetAllProductsResponse(res), nil
}

func toGetAllProductsResponse(products []*models.Product) *pb.GetAllProductsResponse {
	var productList []*pb.Product
	for _, product := range products {
		productList = append(productList, &pb.Product{
			Id:   product.Id.Hex(),
			Name: product.Name,
		})
	}

	return &pb.GetAllProductsResponse{
		Products: productList,
	}
}
