package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/dao"
	"github.com/HollyEllmo/my-first-go-project/pkg/logging"
	pb_prod_products "github.com/HollyEllmo/my-proto-repo/gen/go/prod_service/products/v1"
	"github.com/google/uuid"
)

type Product struct {
	ID            string
	Name          string
	Description   string
	ImageID       *string
	Price         uint64
	CurrencyID    uint32
	Rating        uint32
	CategoryID    uint32
	Specification map[string]interface{}
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}

func NewProduct(ps *dao.ProductStorage) *Product {
	var imageID *string
	if ps.ImageID.Valid {
		imageID = &ps.ImageID.String
	}

	var updatedAt *time.Time
	if ps.UpdatedAt.Valid {
		parsedTime, err := time.Parse(time.RFC3339, ps.UpdatedAt.String)
		if err != nil {
			return nil
		}
		updatedAt = &parsedTime
	}

	createdAt := time.Now()
	if ps.CreatedAt.Valid {
		if parsed, err := time.Parse(time.RFC3339, ps.CreatedAt.String); err == nil {
			createdAt = parsed
		}
	}

	return &Product{
		ID:            ps.ID,
		Name:          ps.Name,
		Description:   ps.Description,
		ImageID:       imageID,
		Price:         ps.Price,
		CurrencyID:    ps.CurrencyID,
		Rating:        ps.Rating,
		CategoryID:    ps.CategoryID,
		Specification: ps.Specification, // Оставляем как map[string]interface{}
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}


func (p *Product) ToMap() (map[string]interface{}, error) {
	updateProductMap := make(map[string]interface{})
	
	// Добавляем только непустые поля
	if p.ID != "" {
		updateProductMap["id"] = p.ID
	}
	if p.Name != "" {
		updateProductMap["name"] = p.Name
	}
	if p.Description != "" {
		updateProductMap["description"] = p.Description
	}
	if p.ImageID != nil {
		updateProductMap["image_id"] = p.ImageID
	}
	if p.Price != 0 {
		updateProductMap["price"] = p.Price
	}
	if p.CurrencyID != 0 {
		updateProductMap["currency_id"] = p.CurrencyID
	}
	if p.Rating != 0 {
		updateProductMap["rating"] = p.Rating
	}
	if p.CategoryID != 0 {
		updateProductMap["category_id"] = p.CategoryID
	}
	if len(p.Specification) > 0 {
		updateProductMap["specification"] = p.Specification
	}
	if !p.CreatedAt.IsZero() {
		updateProductMap["created_at"] = p.CreatedAt
	}
	if p.UpdatedAt != nil {
		updateProductMap["updated_at"] = *p.UpdatedAt
	}
	
	return updateProductMap, nil
}

func NewProductFromPB(productPB *pb_prod_products.CreateProductRequest) (*Product, error) {
	spec, err := parseSpecification(productPB.Specification)
	if err != nil {
		return nil, err
	}
	
	return &Product{
		ID: 		   uuid.New().String(),
		Name:          productPB.GetName(),
		Description:   productPB.GetDescription(),
		ImageID:       productPB.ImageId,
		Price:         productPB.GetPrice(),
		CurrencyID:    productPB.GetCurrencyId(),
		Rating:        productPB.GetRating(),
		CategoryID:    productPB.GetCategoryId(),
		Specification: spec,
		CreatedAt:     time.Now().UTC(),
	}, nil
}



func parseSpecification(specFromPB string) (spec map[string]interface{}, unmarshalErr error) {
	if specFromPB == "" {
		return nil, errors.New("specification is empty")
   }

   if unmarshalErr = json.Unmarshal([]byte(specFromPB), &spec); unmarshalErr == nil {
       return spec, nil
   }

   logging.GetLogger().
       WithError(unmarshalErr).
       WithField("specification", specFromPB).
       Warnf("Failed to unmarshal product specification")

	 return nil, fmt.Errorf("failed to unmarshal product specification: %s", specFromPB)
}

func (p *Product) UpdateFromPB(productPB *pb_prod_products.UpdateProductRequest) {
	if productPB.Name != nil {
		p.Name = productPB.GetName()
	}
	if productPB.Description != nil {
		p.Description = productPB.GetDescription()
	}
	if productPB.ImageId != nil {
		p.ImageID = productPB.ImageId
	}
	if productPB.Price != nil {
		p.Price = productPB.GetPrice()
	}
	if productPB.CurrencyId != nil {
		p.CurrencyID = productPB.GetCurrencyId()
	}
	if productPB.Rating != nil {
		p.Rating = productPB.GetRating()
	}
	if productPB.CategoryId != nil {
		p.CategoryID = productPB.GetCategoryId()
	}

	var spec map[string]interface{}
	if productPB.Specification != nil {
		if err := json.Unmarshal([]byte(productPB.GetSpecification()), &spec); err == nil {
			p.Specification = spec
		} else {
			logging.GetLogger().
			WithError(err).
			WithField("specification", productPB.GetSpecification()).
			Warnf("Failed to unmarshal product specification")
		}
	}
}

func (p Product) ToProto() *pb_prod_products.Product {
	var updatedAt int64
	if p.UpdatedAt != nil {
		updatedAt = p.UpdatedAt.UnixMilli()
	}

	specBytes, err := json.Marshal(p.Specification)
	if err != nil {
		logging.GetLogger().Warnf("Failed to marshal product specification: %v", err)
		logging.GetLogger().Traceln(p.Specification)
		specBytes = []byte("{}") // Default to empty JSON object on error
	}

	return  &pb_prod_products.Product{
		Id:            p.ID,
		Name:          p.Name,
		Description:   p.Description,
		ImageId:       p.ImageID,
		Price:         p.Price,
		CurrencyId:    p.CurrencyID,
		Rating:        p.Rating,
		CategoryId:    p.CategoryID,
		Specification: string(specBytes),
		UpdatedAt:     updatedAt,
		CreatedAt:     p.CreatedAt.UnixMilli(),
	}
}