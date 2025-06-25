package product

import (
	"context"

	"github.com/HollyEllmo/my-first-go-project/internal/controller/dto"
	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/model"
	"github.com/HollyEllmo/my-first-go-project/pkg/logging"
	pb_prod_products "github.com/HollyEllmo/my-proto-repo/gen/go/prod_service/products/v1"
)



func (s *Server) AllProducts(ctx context.Context, request *pb_prod_products.AllProductsRequest) (*pb_prod_products.AllProductsResponse, error) {
	logging.GetLogger().Warningf("ITS IS ALIVE !!!")
	sort := model.ProductsSort(request)
	filter := model.ProductsFilter(request)

	all, err := s.policy.All(ctx, filter, sort)
	if err != nil {
		return  nil, err
	}

	pbProducts := make([]*pb_prod_products.Product, len(all))
	for _, p := range all {
		pbProducts = append(pbProducts, p.ToProto())
	}

	return &pb_prod_products.AllProductsResponse{
		Product: pbProducts,
	}, nil
}

func (s *Server) CreateProduct(ctx context.Context, req *pb_prod_products.CreateProductRequest) (*pb_prod_products.CreateProductResponse, error) {
	d := dto.NewCreateProductDTOFromPB(req)

	product, err := s.policy.CreateProduct(ctx, d)
	if err != nil {
		return  nil, err
	}

	return &pb_prod_products.CreateProductResponse{
		Product: product.ToProto(),
	}, nil
}