package impl

import (
	"Hermes/internal/errs"
	mockLogger "Hermes/internal/logger/mocks"
	"Hermes/internal/models"
	mockStorage "Hermes/internal/repository/mocks"
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNewService(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := mockLogger.NewMockLogger(controller)
	mockStorage := mockStorage.NewMockStorage(controller)

	svc := NewService(mockLogger, mockStorage)

	require.NotNil(t, svc)
	require.Equal(t, mockLogger, svc.logger)
	require.Equal(t, mockStorage, svc.storage)

}

func TestValidateComment(t *testing.T) {

	validComment := models.Comment{Content: "hello", Author: "aboba"}

	t.Run("valid comment", func(t *testing.T) {
		err := validateComment(validComment)
		require.NoError(t, err)
	})

	t.Run("empty content", func(t *testing.T) {
		c := validComment
		c.Content = "   "
		err := validateComment(c)
		require.ErrorIs(t, err, errs.ErrEmptyContent)
	})

	t.Run("empty author", func(t *testing.T) {
		c := validComment
		c.Author = "   "
		err := validateComment(c)
		require.ErrorIs(t, err, errs.ErrEmptyAuthor)
	})

}

func TestService_CreateComment(t *testing.T) {

	ctx := context.Background()
	controller := gomock.NewController(t)

	defer controller.Finish()

	mockLogger := mockLogger.NewMockLogger(controller)
	mockStorage := mockStorage.NewMockStorage(controller)

	svc := &Service{logger: mockLogger, storage: mockStorage}
	comment := models.Comment{Content: "hello", Author: "user"}

	t.Run("validateComment error", func(t *testing.T) {
		invalid := comment
		invalid.Content = ""
		id, err := svc.CreateComment(ctx, invalid)
		require.Equal(t, int64(0), id)
		require.Error(t, err)
	})

	t.Run("storage.CreateComment succeeds", func(t *testing.T) {
		expectedID := int64(123)
		mockStorage.EXPECT().CreateComment(ctx, comment).Return(expectedID, nil)
		id, err := svc.CreateComment(ctx, comment)
		require.NoError(t, err)
		require.Equal(t, expectedID, id)
	})

	t.Run("storage.CreateComment foreign key violation", func(t *testing.T) {
		pgErr := &pgconn.PgError{Code: "23503"}
		mockStorage.EXPECT().CreateComment(ctx, comment).Return(int64(0), pgErr)
		id, err := svc.CreateComment(ctx, comment)
		require.Equal(t, int64(0), id)
		require.ErrorIs(t, err, errs.ErrParentNotFound)
	})

	t.Run("storage.CreateComment generic error", func(t *testing.T) {
		dbErr := errors.New("db down")
		mockStorage.EXPECT().CreateComment(ctx, comment).Return(int64(0), dbErr)
		mockLogger.EXPECT().LogError("service — failed to create comment", dbErr, "id", int64(0), "layer", "service.impl")
		id, err := svc.CreateComment(ctx, comment)
		require.Equal(t, int64(0), id)
		require.EqualError(t, err, "db down")
	})

}

func TestService_DeleteComment(t *testing.T) {

	ctx := context.Background()
	commentID := int64(123)

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := mockLogger.NewMockLogger(controller)
	mockStorage := mockStorage.NewMockStorage(controller)

	svc := &Service{logger: mockLogger, storage: mockStorage}

	t.Run("storage.DeleteComment succeeds", func(t *testing.T) {
		mockStorage.EXPECT().DeleteComment(ctx, commentID).Return(nil)
		err := svc.DeleteComment(ctx, commentID)
		require.NoError(t, err)
	})

	t.Run("storage.DeleteComment ErrCommentNotFound", func(t *testing.T) {
		mockStorage.EXPECT().DeleteComment(ctx, commentID).Return(errs.ErrCommentNotFound)
		err := svc.DeleteComment(ctx, commentID)
		require.ErrorIs(t, err, errs.ErrCommentNotFound)
	})

	t.Run("storage.DeleteComment generic error", func(t *testing.T) {
		dbErr := errors.New("db down")
		mockStorage.EXPECT().DeleteComment(ctx, commentID).Return(dbErr)
		mockLogger.EXPECT().LogError("service — failed to delete comment", dbErr, "id", commentID, "layer", "service.impl")
		err := svc.DeleteComment(ctx, commentID)
		require.EqualError(t, err, "db down")
	})
}

func TestService_GetComments(t *testing.T) {

	ctx := context.Background()
	params := models.QueryParams{}

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := mockLogger.NewMockLogger(controller)
	mockStorage := mockStorage.NewMockStorage(controller)

	svc := &Service{logger: mockLogger, storage: mockStorage}

	t.Run("GetRootComments fails", func(t *testing.T) {
		dbErr := errors.New("db down")
		mockStorage.EXPECT().GetRootComments(ctx, params).Return(nil, dbErr)
		mockLogger.EXPECT().LogError("service — failed to get root comments", dbErr, "layer", "service.impl")
		comments, err := svc.GetComments(ctx, params)
		require.Nil(t, comments)
		require.EqualError(t, err, "db down")
	})

	t.Run("GetCommentTree fails for a root", func(t *testing.T) {
		roots := []models.Comment{{ID: 1}}
		mockStorage.EXPECT().GetRootComments(ctx, params).Return(roots, nil)
		dbErr := errors.New("db down")
		mockStorage.EXPECT().GetCommentTree(ctx, int64(1)).Return(nil, dbErr)
		mockLogger.EXPECT().LogError("service — failed to get comment tree", dbErr, "layer", "service.impl")
		comments, err := svc.GetComments(ctx, params)
		require.Nil(t, comments)
		require.EqualError(t, err, "db down")
	})

	t.Run("success with no roots", func(t *testing.T) {
		mockStorage.EXPECT().GetRootComments(ctx, params).Return([]models.Comment{}, nil)
		comments, err := svc.GetComments(ctx, params)
		require.NoError(t, err)
		require.Empty(t, comments)
	})

	t.Run("success with roots and trees", func(t *testing.T) {

		roots := []models.Comment{{ID: 1}, {ID: 2}}
		flat1 := []models.Comment{{ID: 1}, {ID: 3, ParentID: ptr(1)}}

		mockStorage.EXPECT().GetRootComments(ctx, params).Return(roots, nil)
		mockStorage.EXPECT().GetCommentTree(ctx, int64(1)).Return(flat1, nil)
		mockStorage.EXPECT().GetCommentTree(ctx, int64(2)).Return([]models.Comment{{ID: 2}}, nil)

		comments, err := svc.GetComments(ctx, params)
		require.NoError(t, err)
		require.Len(t, comments, 2)
		require.Equal(t, int64(1), comments[0].ID)
		require.Len(t, comments[0].Children, 1)
		require.Equal(t, int64(3), comments[0].Children[0].ID)
		require.Equal(t, int64(2), comments[1].ID)
		require.Empty(t, comments[1].Children)

	})

}

func TestBuildTree(t *testing.T) {
	t.Run("empty comments", func(t *testing.T) {
		result := buildTree([]models.Comment{})
		require.Empty(t, result)
	})

	t.Run("single root", func(t *testing.T) {
		comments := []models.Comment{{ID: 1}}
		result := buildTree(comments)
		require.Len(t, result, 1)
		require.Equal(t, int64(1), result[0].ID)
		require.Empty(t, result[0].Children)
	})

	t.Run("root with child", func(t *testing.T) {
		comments := []models.Comment{{ID: 1}, {ID: 2, ParentID: ptr(1)}}
		result := buildTree(comments)
		require.Len(t, result, 1)
		require.Equal(t, int64(1), result[0].ID)
		require.Len(t, result[0].Children, 1)
		require.Equal(t, int64(2), result[0].Children[0].ID)
	})

	t.Run("multiple roots", func(t *testing.T) {
		comments := []models.Comment{{ID: 1}, {ID: 2}}
		result := buildTree(comments)
		require.Len(t, result, 2)
	})

	t.Run("nested children", func(t *testing.T) {
		comments := []models.Comment{{ID: 1}, {ID: 2, ParentID: ptr(1)}, {ID: 3, ParentID: ptr(2)}}
		result := buildTree(comments)
		require.Len(t, result, 1)
		require.Len(t, result[0].Children, 1)
		require.Len(t, result[0].Children[0].Children, 1)
		require.Equal(t, int64(3), result[0].Children[0].Children[0].ID)
	})

	t.Run("orphan child treated as root", func(t *testing.T) {
		comments := []models.Comment{{ID: 2, ParentID: ptr(999)}}
		result := buildTree(comments)
		require.Len(t, result, 1)
		require.Equal(t, int64(2), result[0].ID)
	})
}

func ptr(i int64) *int64 {
	return &i
}
