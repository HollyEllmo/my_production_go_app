package dao

import (
	"database/sql"
	"time"

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
		CreatedAt:     time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
	}
}
 
type UpdateProductStorageDTO struct {
	Name          *string `mapstructure:"name, omitempty"`
	Description   *string `mapstructure:"description, omitempty"`
	ImageID       *string `mapstructure:"image_id, omitempty"`
	Price         *uint64 `mapstructure:"price, omitempty"`
	CurrencyID    *uint32 `mapstructure:"currency_id, omitempty"`
	Rating        *uint32 `mapstructure:"rating, omitempty"`
	CategoryID    *uint32 `mapstructure:"category_id, omitempty"`
	Specification map[string]interface{} `mapstructure:"specification, omitempty"`
	UpdatedAt     string `mapstructure:"updated_at, omitempty"`
}

func NewUpdateProductStorageDTO(dto *dto.UpdateProductDTO) *UpdateProductStorageDTO {
	
	return &UpdateProductStorageDTO{
		Name:          dto.Name,
		Description:   dto.Description,
		ImageID:       dto.ImageID,
		Price:         dto.Price,
		CurrencyID:    dto.CurrencyID,
		Rating:        dto.Rating,
		CategoryID:    dto.CategoryID,
		Specification: dto.Specification,
		UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
	}
}