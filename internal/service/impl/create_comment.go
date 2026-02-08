package impl

import (
	"Hermes/internal/errs"
	"Hermes/internal/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Service) CreateComment(ctx context.Context, comment models.Comment) (int64, error) {

	if err := validateComment(comment); err != nil {
		return 0, err
	}

	id, err := s.storage.CreateComment(ctx, comment)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" { // foreign key violation
			return 0, errs.ErrParentNotFound
		}
		s.logger.LogError("service — failed to create comment", err, "id", id, "layer", "service.impl")
		return 0, err
	}

	return id, nil

}
