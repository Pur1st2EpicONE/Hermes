package service

import (
	"Hermes/internal/logger"
	"Hermes/internal/repository"
	"Hermes/internal/service/impl"
)

type Service interface {
}

func NewService(logger logger.Logger, storage repository.Storage) Service {
	return impl.NewService(logger, storage)
}
