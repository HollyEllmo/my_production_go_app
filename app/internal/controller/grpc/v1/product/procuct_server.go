package product

import (
	"context"

	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/model"
	"github.com/HollyEllmo/my-first-go-project/pkg/logging"
	pb_prod_products "github.com/HollyEllmo/my-proto-repo/gen/go/prod_service/products/v1"
)


func (s *Server) AllProducts(ctx context.Context, req *pb_prod_products.AllProductsRequest) (*pb_prod_products.AllProductsResponse, error) {
	logging.GetLogger().Warningf("ITS IS ALIVE !!!")
	sort := model.ProductsSort(req)
	filter := model.ProductsFilter(req)

	all, err := s.policy.All(ctx, filter, sort)
	if err != nil {
		return  nil, err
	}

	pbProducts := make([]*pb_prod_products.Product, len(all))
	for i, p := range all {
		pbProducts[i] = p.ToProto()
	}

	return &pb_prod_products.AllProductsResponse{
		Product: pbProducts,
	}, nil
}

func (s *Server) ProductByID(ctx context.Context, req *pb_prod_products.ProductByIDRequest) (*pb_prod_products.ProductByIDResponse, error) {
 one, err := s.policy.One(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb_prod_products.ProductByIDResponse{
		Product: one.ToProto(),
	}, nil
}

func (s *Server) UpdateProduct(ctx context.Context, req *pb_prod_products.UpdateProductRequest) (*pb_prod_products.UpdateProductResponse, error) {
	product, err := s.policy.One(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	
	product.UpdateFromPB(req)
	
	err = s.policy.Update(ctx, product)
	if err != nil {
		return nil, err
	}

	return &pb_prod_products.UpdateProductResponse{}, nil
}

func (s *Server) DeleteProduct(ctx context.Context, req *pb_prod_products.DeleteProductRequest) (*pb_prod_products.DeleteProductResponse, error) {
	err := s.policy.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb_prod_products.DeleteProductResponse{}, nil
}

func (s *Server) CreateProduct(ctx context.Context, req *pb_prod_products.CreateProductRequest) (*pb_prod_products.CreateProductResponse, error) {
    pb, err := model.NewProductFromPB(req)

	if err != nil {
		logging.WithError(ctx, err)
		return nil, err
	}

	product, err := s.policy.CreateProduct(ctx, pb)
	if err != nil {
		return nil, err
	}

	return &pb_prod_products.CreateProductResponse{
		Product: product.ToProto(),
	}, nil
}