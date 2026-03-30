// Package service defines the business logic layer of the application.
// It orchestrates operations between repository (data layer) and handlers (API layer),
// enforces validation, and maps domain-specific errors.
package service

import (
	"Hermes/internal/logger"
	"Hermes/internal/models"
	"Hermes/internal/repository"
	"Hermes/internal/service/impl"
	"context"
)

// Service defines the interface for comment-related operations.
// It abstracts the underlying implementation and provides
// methods for creating, retrieving, and deleting comments.
type Service interface {
	CreateComment(ctx context.Context, comment models.Comment) (int64, error)                  // Validate and store a new comment, return its ID
	GetComments(ctx context.Context, queryParams models.QueryParams) ([]models.Comment, error) // Retrieve comments with hierarchical tree structure based on query parameters
	DeleteComment(ctx context.Context, id int64) error                                         // Remove a comment by ID, return error if not found
}

// NewService constructs a Service implementation by injecting
// logger and repository.Storage dependencies.
func NewService(logger logger.Logger, storage repository.Storage) Service {
	return impl.NewService(logger, storage)
}
