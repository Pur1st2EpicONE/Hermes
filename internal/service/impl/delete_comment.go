package impl

import (
	"Hermes/internal/errs"
	"context"
	"errors"
)

// DeleteComment removes a comment by ID via the repository layer.
// If the comment does not exist, it returns ErrCommentNotFound.
// Other errors are logged and returned as-is.
func (s *Service) DeleteComment(ctx context.Context, id int64) error {
	if err := s.storage.DeleteComment(ctx, id); err != nil {
		if errors.Is(err, errs.ErrCommentNotFound) {
			return err
		}
		s.logger.LogError("service — failed to delete comment", err, "id", id, "layer", "service.impl")
		return err
	}
	return nil
}
