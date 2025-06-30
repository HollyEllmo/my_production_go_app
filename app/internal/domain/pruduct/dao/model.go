package dao

import (
	"database/sql"
)

type ProductStorage struct {
	ID            string `mapstructure:"id, omitempty"`
	Name          string `mapstructure:"name, omitempty"`
	Description   string `mapstructure:"description, omitempty"`
	ImageID       sql.NullString `mapstructure:"image_id, omitempty"`
	Price         uint64 `mapstructure:"price, omitempty"`
	CurrencyID    uint32 `mapstructure:"currency_id, omitempty"`
	Rating        uint32 `mapstructure:"rating, omitempty"`
	CategoryID    uint32 `mapstructure:"category_id, omitempty"`
	Specification map[string]interface{} `mapstructure:"specification, omitempty"`
	CreatedAt     sql.NullString `mapstructure:"created_at, omitempty"`
	UpdatedAt     sql.NullString `mapstructure:"updated_at, omitempty"`
}

