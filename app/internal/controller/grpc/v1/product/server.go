package product

import (
	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/policy"
	pb_prod_products "github.com/HollyEllmo/my-proto-repo/gen/go/prod_service/products/v1"
)

type Server struct {
	policy *policy.ProductPolicy
	pb_prod_products.UnimplementedProductServiceServer
	
}

func NewServer(policy *policy.ProductPolicy,srv pb_prod_products.UnimplementedProductServiceServer) *Server {
	return &Server{
		policy: policy,
		UnimplementedProductServiceServer: srv,
		
	}
}