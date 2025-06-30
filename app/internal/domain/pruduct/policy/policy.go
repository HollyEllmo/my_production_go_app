package policy

import (
	"context"

	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/model"
	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/service"
	"github.com/HollyEllmo/my-first-go-project/pkg/api/filter"
	"github.com/HollyEllmo/my-first-go-project/pkg/api/sort"
	"github.com/HollyEllmo/my-first-go-project/pkg/errors"
)

type productService interface {
	All(ctx context.Context, filtering filter.Filterable, sorting sort.Sortable) ([]*model.Product, error)
	Create(ctx context.Context, product *model.Product) (*model.Product, error)
	One(ctx context.Context, id string) (*model.Product, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, product *model.Product) error
}

type ProductPolicy struct {
	productService *service.ProductService
}

func NewProductPolicy(productService *service.ProductService) *ProductPolicy {
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

func (p *ProductPolicy) CreateProduct(ctx context.Context, d *model.Product) (*model.Product, error) {
	return p.productService.Create(ctx, d)
}

func (p *ProductPolicy) One(ctx context.Context, id string) (*model.Product, error) {
	product, err := p.productService.One(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "productService.One")
	}

	return product, nil
}

func (p *ProductPolicy) Delete(ctx context.Context, id string) error {
   return p.productService.Delete(ctx, id)
}

func (p *ProductPolicy) Update(ctx context.Context, product *model.Product) error {
   return p.productService.Update(ctx, product)
}

