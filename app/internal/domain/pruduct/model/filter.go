package model

import (
	"github.com/HollyEllmo/my-first-go-project/pkg/api/filter"
	"github.com/HollyEllmo/my-first-go-project/pkg/api/sort"
	pb_prod_products "github.com/HollyEllmo/my_proto_repo/gen/go/prod_service/products/v1"
)

const (
	nameFilterField        = "name"
	descriptionFilterField = "description"
	priceFilterField        = "price"
	ratingFilterField      = "rating"
	categoryIDFilterField = "category_id"
)

func productsFilterFields() map[string]string {
	return map[string]string{
		nameFilterField: filter.DataTypeStr,
		descriptionFilterField: filter.DataTypeStr,
		priceFilterField:      filter.DataTypeStr,
		ratingFilterField:      filter.DataTypeInt,
		categoryIDFilterField:  filter.DataTypeStr,
	}
}

func ProductsSort(req *pb_prod_products.AllProductsRequest) sort.Sortable {
	field := req.GetSort().GetField()
	return sort.NewOptions(field)
}

// func ProductsFilter(req *pb_prod_products.AllProductsRequest) filter.Filterable {
// 	options := filter.NewOptions(
// 		req.GetPagination().GetLimit(),
// 		req.GetPagination().GetOffset(),
// 		productsFilterFields(),
// 	)
// 	if req != nil {
// 		return options
// 	}

// 	Name := req.GetName()
// 	Rating := req.GetRating()
// 	Description := req.GetDescription()
// 	Price := req.GetPrice()
// 	CategoryId := req.GetCategoryId()
// }