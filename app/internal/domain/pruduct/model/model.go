package model

import (
	"encoding/json"
	"time"

	"github.com/HollyEllmo/my-first-go-project/pkg/logging"
	pb_prod_products "github.com/HollyEllmo/my-proto-repo/gen/go/prod_service/products/v1"
)

// NewProduct создает новый экземпляр Product
func NewProduct(
	id, name, description string,
	imageID *string,
	price uint64,
	currencyID, rating uint32,
	categoryID uint32, specification string,
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
	Price         uint64 
	CurrencyID    uint32 
	Rating        uint32 
	CategoryID    uint32 // Changed to uint32 for consistency with proto definition
	Specification string 
	CreatedAt     time.Time 
	UpdatedAt     *time.Time 
}

func (p Product) ToProto() *pb_prod_products.Product {
	var updatedAt int64
	if p.UpdatedAt != nil {
		updatedAt = p.UpdatedAt.UnixMilli()
	}

	specBytes, err := json.Marshal(p.Specification)
	if err != nil {
		logging.GetLogger().Warnf("Failed to marshal product specification: %v", err)
		logging.GetLogger().Tracef(p.Specification)
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