package service

import (
	"context"

	"github.com/HollyEllmo/my-first-go-project/internal/controller/grpc/v1/product"
	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/model"
	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/storage"
)

type repository interface {
	All(ctx context.Context) ([]storage.Product, error)
	Create(ctx context.Context, dto storage.CreateProductDTO) (storage.Product, error)
}

type Service struct {
	repository repository
}

func (s *Service) All(ctx context.Context) ([]model.Product, error) {
	// Implementation for fetching all products
	return nil, nil
}

func (s *Service) Create(ctx context.Context, dto product.CreateProductDTO) (model.Product, error) {
	// Implementation for creating a new product
	return model.Product{}, nil
}