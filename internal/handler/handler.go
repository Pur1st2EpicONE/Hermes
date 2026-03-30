// Package handler provides HTTP handlers for the Hermes application.
// It sets up routes, static file serving, and HTML templates for the web interface.
package handler

import (
	v1 "Hermes/internal/handler/v1"
	"Hermes/internal/service"
	"net/http"
	"text/template"

	"github.com/wb-go/wbf/ginext"
)

const templatePath = "web/templates/index.html"

// NewHandler creates and returns an http.Handler with API routes,
// static file serving and the web frontend at the root path.
func NewHandler(service service.Service) http.Handler {

	handler := ginext.New("")

	handler.Use(ginext.Recovery())
	handler.Static("/static", "./web/static")

	apiV1 := handler.Group("/api/v1")
	handlerV1 := v1.NewHandler(service)

	apiV1.POST("/comments", handlerV1.CreateComment)
	apiV1.GET("/comments", handlerV1.GetComments)
	apiV1.DELETE("/comments/:id", handlerV1.DeleteComment)

	handler.GET("/", homePage(template.Must(template.ParseFiles(templatePath))))

	return handler

}

// homePage renders the main HTML page.
func homePage(t *template.Template) ginext.HandlerFunc {
	return func(c *ginext.Context) {
		if err := t.Execute(c.Writer, nil); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, ginext.H{"error": "Failed to render page"})
		}
	}
}
