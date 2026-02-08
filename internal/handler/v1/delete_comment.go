package v1

import (
	"github.com/wb-go/wbf/ginext"
)

const deleted = "deleted"

func (h *Handler) DeleteComment(c *ginext.Context) {

	id, err := parseQuery(c)
	if err != nil {
		respondError(c, err)
		return
	}

	err = h.service.DeleteComment(c.Request.Context(), *id.ParentID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, deleted)

}
