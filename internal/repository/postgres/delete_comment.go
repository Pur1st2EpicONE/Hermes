package postgres

import (
	"Hermes/internal/errs"
	"context"
	"fmt"

	"github.com/wb-go/wbf/retry"
)

// DeleteComment removes a comment by ID.
// It returns ErrCommentNotFound if no rows were affected.
// Query execution is retried according to configured strategy.
func (s *Storage) DeleteComment(ctx context.Context, id int64) error {

	row, err := s.db.ExecWithRetry(ctx, retry.Strategy(s.config.QueryRetryStrategy), `
	
	DELETE FROM comments 
	WHERE id = $1`,

		id)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	rows, err := row.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get number of affected rows: %w", err)
	}

	if rows == 0 {
		return errs.ErrCommentNotFound
	}

	return nil

}
