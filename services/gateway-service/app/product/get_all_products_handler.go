package product

import (
	"context"

	pb "github.com/kmlcnclk/kc-oms/common/api/product"
	"github.com/kmlcnclk/kc-oms/common/pkg/discovery"
	"github.com/kmlcnclk/kc-oms/common/pkg/handler"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type GetAllProductsHandler struct {
	registry discovery.Registry
	tp       trace.TracerProvider
}

func NewGetAllProductsHandler(registry discovery.Registry, tp trace.TracerProvider) *GetAllProductsHandler {
	return &GetAllProductsHandler{
		registry: registry,
		tp:       tp,
	}

}

func (h *GetAllProductsHandler) Handle(ctx context.Context, req *pb.GetAllProductsRequest) (*handler.SuccessResponse, error) {
	zap.L().Info("Fetching all products")

	conn, err := discovery.ServiceConnection(ctx, "products", h.registry, h.tp)
	if err != nil {
		zap.L().Error("Failed to dial server", zap.Error(err))
	}
	defer conn.Close()

	res, err := pb.NewProductServiceClient(conn).GetAllProducts(ctx, req)

	if err != nil {
		zap.L().Error("Failed to create order", zap.Error(err))
		return nil, err
	}

	zap.L().Info("All products successfully returned!")

	return &handler.SuccessResponse{
		Data: res,
	}, nil
}
