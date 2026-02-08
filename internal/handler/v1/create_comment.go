package v1

import (
	"Hermes/internal/errs"
	"Hermes/internal/models"

	"github.com/wb-go/wbf/ginext"
)

func (h *Handler) CreateComment(c *ginext.Context) {

	var request CreateCommentV1

	if err := c.ShouldBindJSON(&request); err != nil {
		respondError(c, errs.ErrInvalidJSON)
		return
	}

	comment := models.Comment{
		ParentID: request.ParentID,
		Content:  request.Content,
		Author:   request.Author,
	}

	id, err := h.service.CreateComment(c.Request.Context(), comment)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, id)

}
