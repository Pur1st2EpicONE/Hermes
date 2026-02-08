package impl

import (
	"Hermes/internal/errs"
	"Hermes/internal/models"
	"strings"
)

func validateComment(comment models.Comment) error {
	if strings.TrimSpace(comment.Content) == "" {
		return errs.ErrEmptyContent
	}
	if strings.TrimSpace(comment.Author) == "" {
		return errs.ErrEmptyAuthor
	}
	return nil
}
