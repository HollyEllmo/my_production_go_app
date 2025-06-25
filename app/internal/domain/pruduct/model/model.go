package model

import (
	"time"

	pb_prod_products "github.com/HollyEllmo/my-proto-repo/gen/go/prod_service/products/v1"
)

// NewProduct создает новый экземпляр Product
func NewProduct(
	id, name, description string,
	imageID *string,
	price string,
	currencyID, rating uint32,
	categoryID, specification string,
	createdAt time.Time,
	updatedAt *time.Time,
) *Product {
	return &Product{
		ID:            id,
		Name:          name,
		Description:   description,
		ImageID:       imageID,
		Price:         price,
		CurrencyID:    currencyID,
		Rating:        rating,
		CategoryID:    categoryID,
		Specification: specification,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

type Product struct {
	ID            string 
	Name          string 
	Description   string 
	ImageID       *string 
	Price         string 
	CurrencyID    uint32 
	Rating        uint32 
	CategoryID    string 
	Specification string 
	CreatedAt     time.Time 
	UpdatedAt     *time.Time 
}

func (p Product) ToProto() *pb_prod_products.Product {
	var imageID string
	if p.ImageID != nil {
		imageID = *p.ImageID
	}

	var updatedAt int64
	if p.UpdatedAt != nil {
		updatedAt = p.UpdatedAt.UnixMilli()
	}

	return  &pb_prod_products.Product{
		Id:            p.ID,
		Name:          p.Name,
		Description:   p.Description,
		ImageId:       imageID,
		Price:         p.Price,
		CurrencyId:    p.CurrencyID,
		Rating:        p.Rating,
		CategoryId:    p.CategoryID,
		Specification: p.Specification,
		UpdatedAt:     updatedAt,
		CreatedAt:     p.CreatedAt.UnixMilli(),
	}
}