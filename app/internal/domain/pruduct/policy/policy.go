package policy

import (
	"context"

	"github.com/HollyEllmo/my-first-go-project/internal/controller/dto"
	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/model"
	"github.com/HollyEllmo/my-first-go-project/pkg/api/filter"
	"github.com/HollyEllmo/my-first-go-project/pkg/api/sort"
	"github.com/HollyEllmo/my-first-go-project/pkg/errors"
)

type productService interface {
	All(ctx context.Context, filtering filter.Filterable, sorting sort.Sortable) ([]*model.Product, error)
	Create(ctx context.Context, dto *dto.CreateProductDTO) (*model.Product, error)
	One(ctx context.Context, id string) (*model.Product, error)
}

type ProductPolicy struct {
	productService productService
}

func NewProductPolicy(productService productService) *ProductPolicy {
	return &ProductPolicy{
		productService: productService,
	}
}

func (p *ProductPolicy) All(ctx context.Context, filtering filter.Filterable, sorting sort.Sortable) ([]*model.Product, error) {
	products, err := p.productService.All(ctx, filtering, sorting)
	if err != nil {
		return nil, errors.Wrap(err, "productService.All")
	}

	return products, nil
}

func (p *ProductPolicy) CreateProduct(ctx context.Context, d *dto.CreateProductDTO) (*model.Product, error) {
	return p.productService.Create(ctx, d)
}

func (p *ProductPolicy) One(ctx context.Context, id string) (*model.Product, error) {
	product, err := p.productService.One(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "productService.One")
	}

	return product, nil
}

