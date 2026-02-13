package postgres

import (
	"Hermes/internal/models"
	"context"
	"database/sql"
	"fmt"

	"github.com/wb-go/wbf/retry"
)

func (s *Storage) GetRootComments(ctx context.Context, params models.QueryParams) ([]models.Comment, error) {

	order := "created_at DESC"
	if params.Sort == "created_at_asc" {
		order = "created_at ASC"
	}

	var rows *sql.Rows
	var err error

	if params.ParentID == nil {

		rows, err = s.db.QueryWithRetry(ctx, retry.Strategy{
			Attempts: s.config.QueryRetryStrategy.Attempts,
			Delay:    s.config.QueryRetryStrategy.Delay,
			Backoff:  s.config.QueryRetryStrategy.Backoff,
		}, `

            SELECT * FROM comments
            WHERE parent_id IS NULL
            ORDER BY `+order+`
            LIMIT $1 OFFSET $2`,

			params.Limit, params.Offset)

	} else {

		rows, err = s.db.QueryWithRetry(ctx, retry.Strategy{
			Attempts: s.config.QueryRetryStrategy.Attempts,
			Delay:    s.config.QueryRetryStrategy.Delay,
			Backoff:  s.config.QueryRetryStrategy.Backoff,
		}, `

            SELECT * FROM comments
            WHERE id = $1`,

			params.ParentID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	defer func() { _ = rows.Close() }()
	var comments []models.Comment

	for rows.Next() {

		var c models.Comment
		if err := rows.Scan(
			&c.ID,
			&c.ParentID,
			&c.Content,
			&c.Author,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		comments = append(comments, c)

	}

	return comments, nil

}
