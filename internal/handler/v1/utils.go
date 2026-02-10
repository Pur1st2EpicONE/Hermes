package v1

import (
	"Hermes/internal/errs"
	"Hermes/internal/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/wb-go/wbf/ginext"
)

const (
	defaultPage  = 1
	defaultLimit = 20
	defaultSort  = "created_at_desc"
	reverseSort  = "created_at_asc"
	maxLimit     = 100
)

func parseQuery(c *ginext.Context) (models.QueryParams, error) {

	queryParams := models.QueryParams{
		Page:  defaultPage,
		Limit: defaultLimit,
		Sort:  defaultSort,
	}

	if val := c.Query("parent"); val != "" {
		parentID, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return models.QueryParams{}, errs.ErrInvalidParentID
		}
		queryParams.ParentID = &parentID
	}

	if val := c.Query("page"); val != "" {
		page, err := strconv.Atoi(val)
		if err != nil || page < 1 {
			return models.QueryParams{}, errs.ErrInvalidPage
		}
		queryParams.Page = page
	}

	if val := c.Query("limit"); val != "" {
		limit, err := strconv.Atoi(val)
		if err != nil || limit < 1 {
			return models.QueryParams{}, errs.ErrInvalidLimit
		}
		if limit > maxLimit {
			limit = maxLimit
		}
		queryParams.Limit = limit
	}

	if val := c.Query("sort"); val != "" {
		switch val {
		case defaultSort, reverseSort:
			queryParams.Sort = val
		default:
			return models.QueryParams{}, errs.ErrInvalidSort
		}
	}

	queryParams.Offset = (queryParams.Page - 1) * queryParams.Limit

	return queryParams, nil

}

func parseParam(c *ginext.Context) (int64, error) {

	idStr := c.Param("id")
	if idStr == "" {
		return 0, errs.ErrEmptyCommentID
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, errs.ErrInvalidCommentID
	}

	return id, nil

}

func respondOK(c *ginext.Context, response any) {
	c.JSON(http.StatusOK, ginext.H{"result": response})
}

func respondError(c *ginext.Context, err error) {
	if err != nil {
		status, msg := mapErrorToStatus(err)
		c.AbortWithStatusJSON(status, ginext.H{"error": msg})
	}
}

func mapErrorToStatus(err error) (int, string) {

	switch {
	case errors.Is(err, errs.ErrInvalidJSON),
		errors.Is(err, errs.ErrEmptyContent),
		errors.Is(err, errs.ErrEmptyAuthor),
		errors.Is(err, errs.ErrInvalidParentID),
		errors.Is(err, errs.ErrInvalidPage),
		errors.Is(err, errs.ErrInvalidLimit),
		errors.Is(err, errs.ErrEmptyCommentID),
		errors.Is(err, errs.ErrInvalidCommentID),
		errors.Is(err, errs.ErrInvalidSort):
		return http.StatusBadRequest, err.Error()

	case errors.Is(err, errs.ErrParentNotFound),
		errors.Is(err, errs.ErrCommentNotFound):
		return http.StatusNotFound, err.Error()

	default:
		return http.StatusInternalServerError, errs.ErrInternal.Error()
	}

}
