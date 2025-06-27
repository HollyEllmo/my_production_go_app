package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/HollyEllmo/my-first-go-project/internal/controller/dto"
	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/dao"
	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/model"
	"github.com/HollyEllmo/my-first-go-project/pkg/api/filter"
	"github.com/HollyEllmo/my-first-go-project/pkg/api/sort"
	"github.com/HollyEllmo/my-first-go-project/pkg/errors"
)

// convertProductStorageToModel конвертирует ProductStorage в модель Product
func convertProductStorageToModel(ps *dao.ProductStorage) *model.Product {
	var imageID *string
	if ps.ImageID.Valid {
		imageID = &ps.ImageID.String
	}

	var updatedAt *time.Time
	if ps.UpdatedAt.Valid {
		if parsed, err := time.Parse(time.RFC3339, ps.UpdatedAt.String); err == nil {
			updatedAt = &parsed
		}
	}

	createdAt := time.Now()
	if ps.CreatedAt.Valid {
		if parsed, err := time.Parse(time.RFC3339, ps.CreatedAt.String); err == nil {
			createdAt = parsed
		}
	}

	// Конвертируем specification из map в string
	specificationStr := ""
	if ps.Specification != nil {
		if specBytes, err := json.Marshal(ps.Specification); err == nil {
			specificationStr = string(specBytes)
		}
	}

	return &model.Product{
		ID:            ps.ID,
		Name:          ps.Name,
		Description:   ps.Description,
		ImageID:       imageID,
		Price:         ps.Price, // Преобразуем uint64 в string
		CurrencyID:    ps.CurrencyID,
		Rating:        ps.Rating,
		CategoryID:    ps.CategoryID,
		Specification: specificationStr,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

type repository interface {
	All(ctx context.Context, filtering filter.Filterable, sorting sort.Sortable) ([]*dao.ProductStorage, error)
	One(ctx context.Context, id string) (*dao.ProductStorage, error)
	Create(ctx context.Context, dto *dao.CreateProductStorageDTO) error
}

type Service struct {
	repository repository
}

func NewProductService(repository repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) All(ctx context.Context, filtering filter.Filterable, sorting sort.Sortable) ([]*model.Product, error) {
	dbProducts, err := s.repository.All(ctx, filtering, sorting)
	if err != nil {
		return nil, errors.Wrap(err, "repository.All")
	}

	// Конвертируем []*dao.ProductStorage в []*model.Product
	products := make([]*model.Product, len(dbProducts))
	for i, ps := range dbProducts {
		products[i] = convertProductStorageToModel(ps)
	}

	return products, nil
}

func (s *Service) Create(ctx context.Context, d *dto.CreateProductDTO) (*model.Product, error) {
	createProductStorageDTO := dao.NewCreateProductStorageDTO(d)
	
	err := s.repository.Create(ctx, createProductStorageDTO)
	if err != nil {
		return nil, err
	}
	one, err := s.repository.One(ctx, createProductStorageDTO.ID)
	if err != nil {
		return nil, err
	}

	// Используем функцию convertProductStorageToModel для конвертации
	product := convertProductStorageToModel(one)
	return product, nil
}

func (s *Service) One(ctx context.Context, id string) (*model.Product, error) {
	one, err := s.repository.One(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "repository.One")
	}

	product := convertProductStorageToModel(one)
	return product, nil
}