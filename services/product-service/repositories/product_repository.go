package repositories

import (
	"context"

	"github.com/kmlcnclk/kc-oms/common/pkg/models"
	"github.com/kmlcnclk/kc-oms/common/pkg/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type ProductRepositoryInterface interface {
	GetAllProducts(ctx context.Context) ([]*models.Product, error)
}

type ProductRepository struct {
	db             *mongodb.MongoDB
	collectionName string
	tracer         trace.Tracer
}

func NewProductRepository(db *mongodb.MongoDB, collectionName string) *ProductRepository {
	tracer := otel.Tracer("products")
	return &ProductRepository{
		db:             db,
		collectionName: collectionName,
		tracer:         tracer,
	}
}

func (r *ProductRepository) GetAllProducts(ctx context.Context) ([]*models.Product, error) {
	ctx, messageSpan := r.tracer.Start(ctx, "ProductRepository.GetAllProducts")
	defer messageSpan.End()

	collection := r.db.GetCollection(r.collectionName)

	findOptions := options.Find()
	filter := bson.M{}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		zap.L().Error("Failed to fetch products", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*models.Product
	for cursor.Next(ctx) {
		var product models.Product
		if err := cursor.Decode(&product); err != nil {
			zap.L().Error("Failed to decode product", zap.Error(err))
			continue
		}
		products = append(products, &product)
	}

	if err := cursor.Err(); err != nil {
		zap.L().Error("Cursor error", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Products successfully returned!", zap.Int("count", len(products)))

	return products, nil
}
