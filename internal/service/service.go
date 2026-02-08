package service

import (
	"Hermes/internal/logger"
	"Hermes/internal/models"
	"Hermes/internal/repository"
	"Hermes/internal/service/impl"
	"context"
)

type Service interface {
	CreateComment(ctx context.Context, comment models.Comment) (int64, error)
	GetComments(ctx context.Context, queryParams models.QueryParams) ([]models.Comment, error)
	DeleteComment(ctx context.Context, id int64) error
}

func NewService(logger logger.Logger, storage repository.Storage) Service {
	return impl.NewService(logger, storage)
}
