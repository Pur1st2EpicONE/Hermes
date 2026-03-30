package postgres

import (
	"Hermes/internal/models"
	"context"
	"fmt"

	"github.com/wb-go/wbf/retry"
)

// GetCommentTree retrieves a full comment subtree starting from rootID.
// It uses a recursive CTE to fetch all nested comments in a single query,
// ordered by creation time (ascending).
func (s *Storage) GetCommentTree(ctx context.Context, rootID int64) ([]models.Comment, error) {

	rows, err := s.db.QueryWithRetry(ctx, retry.Strategy(s.config.QueryRetryStrategy), `

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
