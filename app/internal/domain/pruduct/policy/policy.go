package policy

import (
	"context"

	"github.com/HollyEllmo/my-first-go-project/internal/controller/grpc/v1/product"
	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/model"
	"github.com/HollyEllmo/my-first-go-project/pkg/api/filter"
	"github.com/HollyEllmo/my-first-go-project/pkg/api/sort"
	"github.com/HollyEllmo/my-first-go-project/pkg/errors"
)

type productService interface {
	All(ctx context.Context, filtering filter.Filterable, sorting sort.Sortable) ([]model.Product, error)
	Create(ctx context.Context, dto product.CreateProductDTO) (model.Product, error)
}

type ProductPolicy struct {
	productService productService
}

func (p *ProductPolicy) All(ctx context.Context, filtering filter.Filterable, sorting sort.Sortable) ([]model.Product, error) {
	products, err := p.productService.All(ctx, filtering, sorting)
	if err != nil {
		return nil, errors.Wrap(err, "productService.All")
	}

	return products, nil
}

