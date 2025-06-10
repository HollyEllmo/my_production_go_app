package handler

import (
	"encoding/json"
	"net/http"
)

// HealthHandler обрабатывает запросы к эндпоинту здоровья
type HealthHandler struct{}

// NewHealthHandler создает новый обработчик здоровья
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health возвращает статус здоровья сервиса
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":  "healthy",
		"service": "my-production-service",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
