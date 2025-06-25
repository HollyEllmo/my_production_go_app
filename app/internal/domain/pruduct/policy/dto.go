package policy

import (
	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/model"
	"github.com/google/uuid"
)

type CreateProductDTO struct {
	ID 		      string
	Name          string 
	Description   string 
	ImageID       *string 
	Price         string 
	CurrencyID    uint32
	Rating        uint32
	CategoryID    string 
	Specification string 
}

func NewCreateProductDTO(dto *model.CreateProductDTO) *CreateProductDTO {
	return &CreateProductDTO{
		ID:            uuid.New().String(),
		Name:          dto.Name,
		Description:   dto.Description,
		ImageID:       dto.ImageID,
		Price:         dto.Price,
		CurrencyID:    dto.CurrencyID,
		Rating:        dto.Rating,
		CategoryID:    dto.CategoryID,
		Specification: dto.Specification,
	}
} 