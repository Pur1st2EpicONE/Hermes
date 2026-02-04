package handler

import (
	v1 "Hermes/internal/handler/v1"
	"Hermes/internal/service"
	"net/http"

	"github.com/wb-go/wbf/ginext"
)

func NewHandler(service service.Service) http.Handler {

	handler := ginext.New("")

	handler.Use(ginext.Recovery())

	apiV1 := handler.Group("/api/v1")
	handlerV1 := v1.NewHandler(service)

	apiV1.POST("/comments", handlerV1.CreateComment)
	apiV1.GET("/comments", handlerV1.GetComments)
	apiV1.DELETE("/comments/:id", handlerV1.DeleteComment)

	return handler

}
