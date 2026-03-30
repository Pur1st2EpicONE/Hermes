// Package v1 contains HTTP handlers for version 1 of the API.
// It defines request handling logic and delegates business operations
// to the service layer.
package v1

import (
	"Hermes/internal/service"
)

// Handler groups all v1 HTTP handlers and holds
// dependencies required to process requests.
type Handler struct {
	service service.Service
}

// NewHandler creates a new v1 Handler instance
// with injected service dependencies.
func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}
