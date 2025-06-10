package metric

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)


const (
	URL = "/api/hearbeat"
)

type Handler struct {}

// TODO fix dependencie on httprouter
func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, URL, h.Heartbeat)
}

// Heartbeat
// @Summary Heartbeat metric
// @Tags Metrics
// @Success 204
// @Failure 400
// @Router /api/hearbeat [get]
func (h *Handler) Heartbeat(w http.ResponseWriter, r *http.Request) {
	log.Print("heartbeat metric called")
	w.WriteHeader(204)
}