package storage

import (
	"database/sql"
	"time"

	"github.com/HollyEllmo/my-first-go-project/internal/controller/dto"
	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/model"
	"github.com/google/uuid"
)

type Product struct {
	ID            string 
	Name          string 
	Description   string 
	ImageID       sql.NullString
	Price         string 
	CurrencyID    uint32 
	Rating        uint32 
	CategoryID    string 
	Specification string 
	CreatedAt     time.Time 
	UpdatedAt     *time.Time 
}

func (p Product) ToModel() model.Product {
	var imageID *string
	
	if p.ImageID.Valid {
		imageID = &p.ImageID.String
	}
	
	return model.Product{
		ID:            p.ID,
		Name:          p.Name,
		Description:   p.Description,
		ImageID:       imageID,
		Price:         p.Price,
		CurrencyID:    p.CurrencyID,
		Rating:        p.Rating,
		CategoryID:    p.CategoryID,
		Specification: p.Specification,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

type CreateProductStorageDTO struct {
	ID			  string
	Name          string 
	Description   string 
	ImageID       sql.NullString
	Price         string 
	CurrencyID    uint32 
	Rating        uint32 
	CategoryID    string 
	Specification string 
	CreatedAt     time.Time 
	UpdatedAt     *time.Time 
}

func NewCreateProductStorageDTO(d *dto.CreateProductDTO) *CreateProductStorageDTO {
	now := time.Now()

    var imageID sql.NullString
	if d.ImageID != nil {
		imageID = sql.NullString{
			String: *d.ImageID,
			Valid: true,
		}
	}

	return  &CreateProductStorageDTO{
		ID: uuid.New().String(),
		Name:          d.Name,
		ImageID: 	   imageID,
		Description:   d.Description,
		Price:         d.Price,
		CurrencyID:    d.CurrencyID,
		Rating:        d.Rating,
		CategoryID:    d.CategoryID,
		Specification: d.Specification,
		CreatedAt:     now,
		UpdatedAt:     &now,
	}
}