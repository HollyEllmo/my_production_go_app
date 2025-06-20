package product

import (
	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/model"
	pb_prod_products "github.com/HollyEllmo/my-proto-repo/gen/go/prod_service/products/v1"
)

// CreateProductDTO is now in the model package
// This is a wrapper function for convenience
func NewCreateProductDTOFromPB(product *pb_prod_products.CreateProductRequest) *model.CreateProductDTO {
   return model.NewCreateProductDTOFromPB(product)
}