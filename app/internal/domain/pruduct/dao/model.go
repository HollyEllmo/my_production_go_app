package dao

import (
	"database/sql"

	"github.com/HollyEllmo/my-first-go-project/internal/controller/dto"
	"github.com/google/uuid"
)

type ProductStorage struct {
	ID            string
	Name          string
	Description   string
	ImageID       sql.NullString
	Price         uint64
	CurrencyID    uint32
	Rating        uint32
	CategoryID    uint32 // Changed to int32 for consistency with proto definition
	Specification map[string]interface{}
	CreatedAt     sql.NullString
	UpdatedAt     sql.NullString
}

type CreateProductStorageDTO struct {
	ID            string
	Name          string
	Description   string
	ImageID       *string
	Price         uint64
	CurrencyID    uint32
	Rating        uint32
	CategoryID    uint32
	Specification map[string]interface{}
	CreatedAt     string
	UpdatedAt     string
}

// NewCreateProductStorageDTO создает DTO для создания продукта в хранилище
func NewCreateProductStorageDTO(dto *dto.CreateProductDTO) *CreateProductStorageDTO {
	return &CreateProductStorageDTO{
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
