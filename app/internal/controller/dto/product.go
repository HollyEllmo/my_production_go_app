package dto

import (
	"encoding/json"

	"github.com/HollyEllmo/my-first-go-project/pkg/logging"
	pb_prod_products "github.com/HollyEllmo/my-proto-repo/gen/go/prod_service/products/v1"
)

type CreateProductDTO struct {
	Name          string 
	Description   string 
	ImageID       *string 
	Price         uint64 
	CurrencyID    uint32
	Rating        uint32
	CategoryID    uint32
	Specification map[string]interface{}
}

func NewCreateProductDTOFromPB(product *pb_prod_products.CreateProductRequest) *CreateProductDTO {
	var spec map[string]interface{}
	err := json.Unmarshal([]byte(product.Specification), &spec)
	if err != nil {
		spec = make(map[string]interface{})
	}

	return &CreateProductDTO{
		Name:          product.GetName(),
		Description:   product.GetDescription(),
		ImageID:       product.ImageId,
		Price:         product.GetPrice(),
		CurrencyID:    product.GetCurrencyId(),
		Rating:        product.GetRating(),
		CategoryID:    product.GetCategoryId(),
		Specification: spec,
	}
}

type UpdateProductDTO struct {
	Name          *string `mapstructure:"name, omitempty"`
	Description   *string `mapstructure:"description, omitempty"`
	ImageID       *string `mapstructure:"image_id, omitempty"`
	Price         *uint64 `mapstructure:"price, omitempty"`
	CurrencyID    *uint32 `mapstructure:"currency_id, omitempty"`
	Rating        *uint32 `mapstructure:"rating, omitempty"`
	CategoryID    *uint32 `mapstructure:"category_id, omitempty"`
	Specification map[string]interface{} `mapstructure:"specification, omitempty"`
}

func NewUpdateProductDTOFromPB(product *pb_prod_products.UpdateProductRequest) *UpdateProductDTO {
	var spec map[string]interface{}
	if product.Specification != nil {
		err := json.Unmarshal([]byte(*product.Specification), &spec)
	  if err != nil {
		logging.GetLogger().Warnf("failed to unmarshal specification: %v", err)
		logging.GetLogger().Traceln(product.Specification)
		spec = make(map[string]interface{})
	  }
	}

	return &UpdateProductDTO{
		Name:          product.Name,
		Description:   product.Description,
		ImageID:       product.ImageId,
		Price:         product.Price,
		CurrencyID:    product.CurrencyId,
		Rating:        product.Rating,
		CategoryID:    product.CategoryId,
		Specification: spec,
	}
}
