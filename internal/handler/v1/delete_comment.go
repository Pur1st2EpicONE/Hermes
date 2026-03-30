package v1

import (
	"github.com/wb-go/wbf/ginext"
)

const deleted = "deleted"

// DeleteComment handles DELETE /comments/:id requests.
// It parses the comment ID from the path, invokes the service layer
// to delete the comment, and returns a confirmation response.
func (h *Handler) DeleteComment(c *ginext.Context) {

	id, err := parseParam(c)
	if err != nil {
		respondError(c, err)
		return
	}

	err = h.service.DeleteComment(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, deleted)

}
