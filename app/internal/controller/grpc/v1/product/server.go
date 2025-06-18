package product

import pb_prod_products "github.com/HollyEllmo/my_proto_repo/gen/go/prod_service/products/v1"

type Server struct {
	pb_prod_products.UnimplementedProductServiceServer
}

func NewServer(srv pb_prod_products.UnimplementedProductServiceServer) *Server {
	return &Server{
		UnimplementedProductServiceServer: srv,
	}
}