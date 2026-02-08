package postgres

import (
	"Hermes/internal/models"
	"context"
	"fmt"

	"github.com/wb-go/wbf/retry"
)

func (s *Storage) CreateComment(ctx context.Context, comment models.Comment) (int64, error) {

	row, err := s.db.QueryRowWithRetry(ctx, retry.Strategy{
		Attempts: s.config.QueryRetryStrategy.Attempts,
		Delay:    s.config.QueryRetryStrategy.Delay,
		Backoff:  s.config.QueryRetryStrategy.Backoff}, `
		
		INSERT INTO comments (parent_id, content, author)
		VALUES ($1, $2, $3)
		RETURNING id`,

		comment.ParentID, comment.Content, comment.Author)
	if err != nil {
		return 0, err
	}

	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to scan row: %w", err)
	}

	return id, nil

}
