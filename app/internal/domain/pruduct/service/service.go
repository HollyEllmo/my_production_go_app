package service

import (
	"context"

	"github.com/HollyEllmo/my-first-go-project/internal/controller/grpc/v1/product"
	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/model"
	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/storage"
	"github.com/HollyEllmo/my-first-go-project/pkg/api/filter"
	"github.com/HollyEllmo/my-first-go-project/pkg/api/sort"
	"github.com/HollyEllmo/my-first-go-project/pkg/errors"
)

type repository interface {
	All(ctx context.Context, filtering filter.Filterable, sorting sort.Sortable) ([]storage.Product, error)
	Create(ctx context.Context, dto storage.CreateProductDTO) (storage.Product, error)
}

type Service struct {
	repository repository
}

func (s *Service) All(ctx context.Context, filtering filter.Filterable, sorting sort.Sortable) ([]model.Product, error) {
	dbProsucts, err := s.repository.All(ctx, filtering, sorting)
	if err != nil {
		return nil, errors.Wrap(err, "repository.All")
	}

	var products []model.Product
	for _, dbP := range dbProsucts {
		products = append(products, dbP.ToModel())
	}
	return products, nil
}

func (s *Service) Create(ctx context.Context, dto product.CreateProductDTO) (model.Product, error) {
	// Implementation for creating a new product
	return model.Product{}, nil
}