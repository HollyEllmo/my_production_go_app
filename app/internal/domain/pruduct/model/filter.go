package model

import (
	"github.com/HollyEllmo/my-first-go-project/internal/controller/grpc/types"
	"github.com/HollyEllmo/my-first-go-project/pkg/api/filter"
	"github.com/HollyEllmo/my-first-go-project/pkg/api/sort"
	"github.com/HollyEllmo/my-first-go-project/pkg/logging"
	pb_prod_products "github.com/HollyEllmo/my-proto-repo/gen/go/prod_service/products/v1"
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

func ProductsFilter(req *pb_prod_products.AllProductsRequest) filter.Filterable {
	options := filter.NewOptions(
		req.GetPagination().GetLimit(),
		req.GetPagination().GetOffset(),
		productsFilterFields(),
	)
	if req != nil {
		return options
	}

	name := req.GetName()
	if name != nil {
		operator := types.StringOperatorFromPB(name.GetOp())
		addFilterField(nameFilterField, name.GetVal(), operator, options)
	}

	rating := req.GetRating()
	if rating != nil {
		operator := types.IntOperatorFromPB(rating.GetOp())
		addFilterField(ratingFilterField, rating.GetVal(), operator, options)
	}
	description := req.GetDescription()
	if description != nil {
		operator := types.StringOperatorFromPB(description.GetOp())
		addFilterField(descriptionFilterField, description.GetVal(), operator, options)
	}
	price := req.GetPrice()
	if price != nil {
		operator := types.IntOperatorFromPB(price.GetOp())
		addFilterField(priceFilterField, price.GetVal(), operator, options)
	}
	categoryId := req.GetCategoryId()
	if categoryId != nil {
		operator := types.IntOperatorFromPB(categoryId.GetOp())
		addFilterField(categoryIDFilterField, categoryId.GetVal(), operator, options)
	}
}

func addFilterField(name, value string, 
	operator filter.Operator,
	options filter.Filterable,
	) {
	err := options.AddField(name, operator, value)
		if err != nil {
			logging.GetLogger().WithError(err).Errorf("failed to add filter field. name=%s, operator=%s, value=%s",
				name, operator, value)
		}
}


