package v1

import (
	"github.com/wb-go/wbf/ginext"
)

func (h *Handler) GetComments(c *ginext.Context) {

	queryParams, err := parseQuery(c)
	if err != nil {
		respondError(c, err)
		return
	}

	comments, err := h.service.GetComments(c.Request.Context(), queryParams)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, comments)

}
