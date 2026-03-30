package impl

import (
	"Hermes/internal/errs"
	"Hermes/internal/models"
	"strings"
)

// validateComment ensures that mandatory fields of a comment are present.
// Returns ErrEmptyContent or ErrEmptyAuthor if validation fails.
func validateComment(comment models.Comment) error {
	if strings.TrimSpace(comment.Content) == "" {
		return errs.ErrEmptyContent
	}
	if strings.TrimSpace(comment.Author) == "" {
		return errs.ErrEmptyAuthor
	}
	return nil
}
