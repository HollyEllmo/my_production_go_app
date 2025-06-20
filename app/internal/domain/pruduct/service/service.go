package service

import (
	"context"

	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/model"
	"github.com/HollyEllmo/my-first-go-project/pkg/api/filter"
	"github.com/HollyEllmo/my-first-go-project/pkg/api/sort"
	"github.com/HollyEllmo/my-first-go-project/pkg/errors"
)

type repository interface {
	All(ctx context.Context, filtering filter.Filterable, sorting sort.Sortable) ([]model.Product, error)
}

type Service struct {
	repository repository
}

func NewProductService(repository repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) All(ctx context.Context, filtering filter.Filterable, sorting sort.Sortable) ([]model.Product, error) {
	products, err := s.repository.All(ctx, filtering, sorting)
	if err != nil {
		return nil, errors.Wrap(err, "repository.All")
	}

	return products, nil
}

func (s *Service) Create(ctx context.Context, dto model.CreateProductDTO) (model.Product, error) {
	// Implementation for creating a new product
	return model.Product{}, nil
}