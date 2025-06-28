package model

import (
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
	Specification string 
}

func NewCreateProductDTOFromPB(product *pb_prod_products.CreateProductRequest) *CreateProductDTO {
	// Преобразуем price из string в uint64
	// price, err := strconv.ParseUint(product.GetPrice(), 10, 64)
	// if err != nil {
	// 	price = 0 // значение по умолчанию в случае ошибки
	// }

	// Преобразуем categoryId из string в uint32
	// categoryID, err := strconv.ParseUint(product.GetCategoryId(), 10, 32)
	// if err != nil {
	// 	categoryID = 0 // значение по умолчанию в случае ошибки
	// }

	return &CreateProductDTO{
		Name:          product.GetName(),
		Description:   product.GetDescription(),
		ImageID:       product.ImageId,
		Price:         product.GetPrice(),
		CurrencyID:    product.GetCurrencyId(),
		Rating:        product.GetRating(),
		CategoryID:    product.GetCategoryId(),
		Specification: product.GetSpecification(),
	}
}
