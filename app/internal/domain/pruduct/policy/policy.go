package policy

import (
	"context"

	"github.com/HollyEllmo/my-first-go-project/internal/controller/grpc/v1/product"
	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/model"
)

type productService interface {
	All(ctx context.Context) ([]model.Product, error)
	Create(ctx context.Context, dto product.CreateProductDTO) (model.Product, error)
}

type ProductPolicy struct {
	productService productService
}