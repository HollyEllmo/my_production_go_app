package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/HollyEllmo/my-first-go-project/internal/controller/dto"
	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/dao"
	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/model"
	"github.com/HollyEllmo/my-first-go-project/pkg/api/filter"
	"github.com/HollyEllmo/my-first-go-project/pkg/api/sort"
	"github.com/HollyEllmo/my-first-go-project/pkg/errors"
)

// convertProductStorageToModel конвертирует ProductStorage в модель Product
func convertProductStorageToModel(ps *dao.ProductStorage) model.Product {
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

	return model.Product{
		ID:            ps.ID,
		Name:          ps.Name,
		Description:   ps.Description,
		ImageID:       imageID,
		Price:         strconv.FormatUint(ps.Price, 10), // Преобразуем uint64 в string
		CurrencyID:    ps.CurrencyID,
		Rating:        ps.Rating,
		CategoryID:    strconv.FormatUint(uint64(ps.CategoryID), 10), // Преобразуем uint32 в string
		Specification: specificationStr,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

type repository interface {
	All(ctx context.Context, filtering filter.Filterable, sorting sort.Sortable) ([]*dao.ProductStorage, error)
	One(ctx context.Context, ID string) (*dao.ProductStorage, error)
	Create(ctx context.Context, m map[string]interface{}) error
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
	productsStorage, err := s.repository.All(ctx, filtering, sorting)
	if err != nil {
		return nil, errors.Wrap(err, "repository.All")
	}

	// Конвертируем []*dao.ProductStorage в []model.Product
	products := make([]model.Product, len(productsStorage))
	for i, ps := range productsStorage {
		products[i] = convertProductStorageToModel(ps)
	}

	return products, nil
}

func (s *Service) Create(ctx context.Context, d *dto.CreateProductDTO) (*model.Product, error) {
	createProductStorageDTO := dao.NewCreateProductStorageDTO(d)
	
	// Конвертируем DTO в map для вызова метода Create
	m := map[string]interface{}{
		"id":            createProductStorageDTO.ID,
		"name":          createProductStorageDTO.Name,
		"description":   createProductStorageDTO.Description,
		"price":         createProductStorageDTO.Price,
		"currency_id":   createProductStorageDTO.CurrencyID,
		"rating":        createProductStorageDTO.Rating,
		"category_id":   createProductStorageDTO.CategoryID,
		"specification": createProductStorageDTO.Specification,
	}
	
	if createProductStorageDTO.ImageID != nil {
		m["image_id"] = *createProductStorageDTO.ImageID
	}
	
	err := s.repository.Create(ctx, m)
	if err != nil {
		return nil, err
	}
	one, err := s.repository.One(ctx, createProductStorageDTO.ID)
	if err != nil {
		return nil, err
	}

	// Используем функцию convertProductStorageToModel для конвертации
	product := convertProductStorageToModel(one)
	return &product, nil
}

func (s *Service) One(ctx context.Context, ID string) (*model.Product, error) {
	productStorage, err := s.repository.One(ctx, ID)
	if err != nil {
		return nil, errors.Wrap(err, "repository.One")
	}

	product := convertProductStorageToModel(productStorage)
	return &product, nil
}