package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MatheusHenrique129/bemax-backend/internal/adapters/constants"
)

type healthHandler struct{}

type HealthHandler interface {
	Ping(w http.ResponseWriter, r *http.Request)
}

func (h healthHandler) Ping(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "pong"}
	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func NewHealthHandler() HealthHandler {
	return &healthHandler{}
}
