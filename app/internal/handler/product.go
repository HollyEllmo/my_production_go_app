package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/HollyEllmo/my-first-go-project/pkg/logging"
	pb_prod_products "github.com/HollyEllmo/my-proto-repo/gen/go/prod_service/products/v1"
	"github.com/julienschmidt/httprouter"
)

type ProductHandler struct {
	client pb_prod_products.ProductServiceClient
}

func NewProductHandler(client pb_prod_products.ProductServiceClient) *ProductHandler {
	return &ProductHandler{
		client: client,
	}
}

// Register регистрирует маршруты для продуктов
func (h *ProductHandler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/api/v1/products", h.GetAllProducts)
	router.HandlerFunc(http.MethodPost, "/api/v1/products", h.CreateProduct)
}

// GetAllProducts получает все продукты через gRPC
func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	if h.client == nil {
		logging.GetLogger(ctx).Errorln("Product service client is not available")
		http.Error(w, "Product service unavailable", http.StatusServiceUnavailable)
		return
	}

	// Создаем контекст с таймаутом
	grpcCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Вызываем gRPC метод
	req := &pb_prod_products.AllProductsRequest{}
	resp, err := h.client.AllProducts(grpcCtx, req)
	if err != nil {
		logging.GetLogger(ctx).WithError(err).Errorln("Failed to get products from gRPC service")
		http.Error(w, "Failed to get products", http.StatusInternalServerError)
		return
	}

	// Возвращаем результат
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp.Product); err != nil {
		logging.GetLogger(ctx).WithError(err).Errorln("Failed to encode response")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	logging.GetLogger(ctx).Infof("Successfully returned %d products", len(resp.Product))
}

// CreateProduct создает новый продукт через gRPC
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	if h.client == nil {
		logging.GetLogger(ctx).Errorln("Product service client is not available")
		http.Error(w, "Product service unavailable", http.StatusServiceUnavailable)
		return
	}

	// Парсим JSON из запроса
	var req pb_prod_products.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logging.GetLogger(ctx).WithError(err).Errorln("Failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Создаем контекст с таймаутом
	grpcCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Вызываем gRPC метод
	resp, err := h.client.CreateProduct(grpcCtx, &req)
	if err != nil {
		logging.GetLogger(ctx).WithError(err).Errorln("Failed to create product via gRPC service")
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	// Возвращаем результат
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp.Product); err != nil {
		logging.GetLogger(ctx).WithError(err).Errorln("Failed to encode response")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	logging.GetLogger(ctx).Infof("Successfully created product with ID: %s", resp.Product.Id)
}
