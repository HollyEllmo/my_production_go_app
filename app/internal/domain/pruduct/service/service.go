package service

import (
	"context"

	"github.com/HollyEllmo/my-first-go-project/internal/controller/dto"
	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/model"
	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/storage"
	"github.com/HollyEllmo/my-first-go-project/pkg/api/filter"
	"github.com/HollyEllmo/my-first-go-project/pkg/api/sort"
	"github.com/HollyEllmo/my-first-go-project/pkg/errors"
)

type repository interface {
	All(ctx context.Context, filtering filter.Filterable, sorting sort.Sortable) ([]model.Product, error)
	One(ctx context.Context, ID string) (*storage.Product, error)
	Create(ctx context.Context, createProductStorageDto *storage.CreateProductStorageDTO) error
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

func (s *Service) Create(ctx context.Context, d *dto.CreateProductDTO) (*model.Product, error) {
	createProductStorageDTO := storage.NewCreateProductStorageDTO(d)
	err := s.repository.Create(ctx, createProductStorageDTO)
	if err != nil {
		return  nil, err
	}
	one, err := s.repository.One(ctx, createProductStorageDTO.ID)
	if err != nil {
		return  nil, err
	}

	// Используем метод ToModel() для конвертации
	product := one.ToModel()
	return &product, nil
}