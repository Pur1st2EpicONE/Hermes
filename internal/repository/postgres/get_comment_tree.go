package postgres

import (
	"Hermes/internal/models"
	"context"
	"fmt"

	"github.com/wb-go/wbf/retry"
)

func (s *Storage) GetCommentTree(ctx context.Context, rootID int64) ([]models.Comment, error) {

	rows, err := s.db.QueryWithRetry(ctx, retry.Strategy{
		Attempts: s.config.QueryRetryStrategy.Attempts,
		Delay:    s.config.QueryRetryStrategy.Delay,
		Backoff:  s.config.QueryRetryStrategy.Backoff,
	}, `

	    WITH RECURSIVE tree AS (
    
		SELECT *
        FROM comments
        WHERE id = $1

        UNION ALL

        SELECT c.*
        FROM comments c
        JOIN tree t ON c.parent_id = t.id
		
		)

        SELECT * FROM tree
        ORDER BY created_at ASC
	
	`, rootID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	defer func() { _ = rows.Close() }()
	var comments []models.Comment

	for rows.Next() {

		var comment models.Comment
		if err := rows.Scan(
			&comment.ID,
			&comment.ParentID,
			&comment.Content,
			&comment.Author,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		comments = append(comments, comment)

	}

	return comments, nil

}
