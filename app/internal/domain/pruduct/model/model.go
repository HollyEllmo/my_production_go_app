package model

import (
	"time"

	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/storage"
	pb_prod_products "github.com/HollyEllmo/my-proto-repo/gen/go/prod_service/products/v1"
)

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

func NewProduct(sp storage.Product) *Product {
	var imageID *string
	if sp.ImageID.Valid {
		imageID = &sp.ImageID.String
	}
	return  &Product{
		ID:            sp.ID,
		Name:          sp.Name,
		Description:   sp.Description,
		ImageID:       imageID,
		Price:         sp.Price,
		CurrencyID:    sp.CurrencyID,
		Rating:        sp.Rating,
		CategoryID:    sp.CategoryID,
		Specification: sp.Specification,
		CreatedAt:     sp.CreatedAt,
		UpdatedAt:     sp.UpdatedAt,
	}
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