package v1

import (
	"Hermes/internal/errs"
	"Hermes/internal/models"
	mockService "Hermes/internal/service/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wb-go/wbf/ginext"
	"go.uber.org/mock/gomock"
)

func setupRouter(handler *Handler) *ginext.Engine {

	r := ginext.New("")

	v1 := r.Group("/api/v1")
	{
		v1.POST("/comments", handler.CreateComment)
		v1.GET("/comments", handler.GetComments)
		v1.DELETE("/comments/:id", handler.DeleteComment)
	}

	return r

}

func TestHandler_CreateComment(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mockService.NewMockService(ctrl)

	h := &Handler{service: mockService}
	router := setupRouter(h)

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/comments", bytes.NewBufferString(`{invalid}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service returns error", func(t *testing.T) {
		body := CreateCommentV1{ParentID: new(int64), Content: "test", Author: "author"}
		*body.ParentID = 999
		b, _ := json.Marshal(body)
		mockService.EXPECT().CreateComment(gomock.Any(), models.Comment{ParentID: body.ParentID, Content: body.Content, Author: body.Author}).Return(int64(0), errs.ErrParentNotFound)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/comments", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("success", func(t *testing.T) {
		body := CreateCommentV1{Content: "test", Author: "author"}
		b, _ := json.Marshal(body)
		mockService.EXPECT().CreateComment(gomock.Any(), models.Comment{Content: body.Content, Author: body.Author}).Return(int64(123), nil)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/comments", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, w.Body.String(), "123")
	})

}

func TestHandler_DeleteComment(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mockService.NewMockService(ctrl)

	h := &Handler{service: mockService}
	router := setupRouter(h)

	t.Run("comment not found", func(t *testing.T) {
		mockService.EXPECT().DeleteComment(gomock.Any(), int64(999)).Return(errs.ErrCommentNotFound)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/comments/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("success delete", func(t *testing.T) {
		mockService.EXPECT().DeleteComment(gomock.Any(), int64(123)).Return(nil)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/comments/123", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, w.Body.String(), "deleted")
	})

}

func TestHandler_GetComments(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mockService.NewMockService(ctrl)

	h := &Handler{service: mockService}
	router := setupRouter(h)

	comments := []models.Comment{{ID: 1, Content: "comment1"}, {ID: 2, Content: "comment2"}}

	t.Run("invalid query", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/comments?sort=invalid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("internal error", func(t *testing.T) {
		qp := models.QueryParams{Page: 1, Limit: 20, Sort: "created_at_desc", Offset: 0}
		mockService.EXPECT().GetComments(gomock.Any(), qp).Return(nil, errors.New("db down :("))
		req := httptest.NewRequest(http.MethodGet, "/api/v1/comments", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("success", func(t *testing.T) {
		qp := models.QueryParams{Page: 1, Limit: 20, Sort: "created_at_desc", Offset: 0}
		mockService.EXPECT().GetComments(gomock.Any(), qp).Return(comments, nil)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/comments", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()
		require.Contains(t, body, `"id":1`)
		require.Contains(t, body, `"content":"comment1"`)
		require.Contains(t, body, `"id":2`)
		require.Contains(t, body, `"content":"comment2"`)
	})

}
