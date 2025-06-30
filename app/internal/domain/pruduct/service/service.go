package service

import (
	"context"

	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/dao"
	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/model"
	"github.com/HollyEllmo/my-first-go-project/pkg/api/filter"
	"github.com/HollyEllmo/my-first-go-project/pkg/api/sort"
	"github.com/HollyEllmo/my-first-go-project/pkg/errors"
)

type repository interface {
	All(ctx context.Context, filtering filter.Filterable, sorting sort.Sortable) ([]*dao.ProductStorage, error)
	One(ctx context.Context, id string) (*dao.ProductStorage, error)
	Create(ctx context.Context, m map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, dm map[string]interface{}) error
}

type ProductService struct {
	repository repository
}

func NewProductService(repository repository) *ProductService {
	return &ProductService{
		repository: repository,
	}
}

func (s *ProductService) All(ctx context.Context, filtering filter.Filterable, sorting sort.Sortable) ([]*model.Product, error) {
	dbProducts, err := s.repository.All(ctx, filtering, sorting)
	if err != nil {
		return nil, errors.Wrap(err, "repository.All")
	}

	// Конвертируем []*dao.ProductStorage в []*model.Product
	products := make([]*model.Product, len(dbProducts))
	for i, ps := range dbProducts {
		products[i] = model.NewProduct(ps)
	}

	return products, nil
}

func (s *ProductService) Create(ctx context.Context, product *model.Product) (*model.Product, error) {
	toMap, err := product.ToMap()
	if err != nil {
		return nil, err
	}

	err = s.repository.Create(ctx, toMap)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) One(ctx context.Context, id string) (*model.Product, error) {
	one, err := s.repository.One(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "repository.One")
	}

	product := model.NewProduct(one)
	return product, nil
}

func (s *ProductService) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}

func (s *ProductService) Update(ctx context.Context, product *model.Product) error {
   toMap, err := product.ToMap()
	if err != nil {
		return  err
	}
	// Обновляем продукт в репозитории
	return s.repository.Update(ctx, product.ID, toMap)
}

// convertProductStorageToModel конвертирует ProductStorage в модель Product
// func convertProductStorageToModel(ps *dao.ProductStorage) *model.Product {
// 	var imageID *string
// 	if ps.ImageID.Valid {
// 		imageID = &ps.ImageID.String
// 	}

// 	var updatedAt *time.Time
// 	if ps.UpdatedAt.Valid {
// 		if parsed, err := time.Parse(time.RFC3339, ps.UpdatedAt.String); err == nil {
// 			updatedAt = &parsed
// 		}
// 	}

// 	createdAt := time.Now()
// 	if ps.CreatedAt.Valid {
// 		if parsed, err := time.Parse(time.RFC3339, ps.CreatedAt.String); err == nil {
// 			createdAt = parsed
// 		}
// 	}

// 	return &model.Product{
// 		ID:            ps.ID,
// 		Name:          ps.Name,
// 		Description:   ps.Description,
// 		ImageID:       imageID,
// 		Price:         ps.Price,
// 		CurrencyID:    ps.CurrencyID,
// 		Rating:        ps.Rating,
// 		CategoryID:    ps.CategoryID,
// 		Specification: ps.Specification, // Оставляем как map[string]interface{}
// 		CreatedAt:     createdAt,
// 		UpdatedAt:     updatedAt,
// 	}
// }