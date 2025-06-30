package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/HollyEllmo/my-first-go-project/pkg/logging"
	pb_prod_products "github.com/HollyEllmo/my-proto-repo/gen/go/prod_service/products/v1"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
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

func (p *Product) ToMap() (map[string]interface{}, error) {
	var updateProductMap map[string]interface{}
	err := mapstructure.Decode(p, &updateProductMap)
	if err != nil {
		return updateProductMap, fmt.Errorf("mapstructure.Decode(product): %w", err)
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
	if specFromPB != "" {
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