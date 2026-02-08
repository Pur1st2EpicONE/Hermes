package postgres

import (
	"Hermes/internal/errs"
	"context"
	"fmt"

	"github.com/wb-go/wbf/retry"
)

func (s *Storage) DeleteComment(ctx context.Context, id int64) error {

	row, err := s.db.ExecWithRetry(ctx, retry.Strategy{
		Attempts: s.config.QueryRetryStrategy.Attempts,
		Delay:    s.config.QueryRetryStrategy.Delay,
		Backoff:  s.config.QueryRetryStrategy.Backoff,
	}, `
	
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
