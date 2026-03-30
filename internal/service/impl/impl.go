// Package impl contains the concrete implementation of the Service interface.
// It handles validation, error mapping, tree-building for nested comments,
// and logging of business-layer errors.
package impl

import (
	"Hermes/internal/logger"
	"Hermes/internal/repository"
)

// Service is the Postgres-backed implementation of the Service interface.
// It wraps a repository.Storage for data access and a logger for error reporting.
type Service struct {
	logger  logger.Logger
	storage repository.Storage
}

// NewService constructs a Service instance with injected logger and storage.
func NewService(logger logger.Logger, storage repository.Storage) *Service {
	return &Service{logger: logger, storage: storage}
}
