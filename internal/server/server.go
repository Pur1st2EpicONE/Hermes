package server

import (
	"Hermes/internal/config"
	"Hermes/internal/logger"
	"Hermes/internal/server/httpserver"
	"net/http"
)

type Server interface {
	Run() error
	Shutdown()
}

func NewServer(logger logger.Logger, config config.Server, handler http.Handler) Server {
	return httpserver.NewServer(logger, config, handler)
}
