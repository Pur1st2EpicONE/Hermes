package repository

import (
	"Hermes/internal/config"
	"Hermes/internal/logger"
	"Hermes/internal/models"
	"Hermes/internal/repository/postgres"
	"context"
	"fmt"

	"github.com/wb-go/wbf/dbpg"
)

type Storage interface {
	Close()
	CreateComment(ctx context.Context, comment models.Comment) (int64, error)
	GetRootComments(ctx context.Context, queryParams models.QueryParams) ([]models.Comment, error)
	GetCommentTree(ctx context.Context, id int64) ([]models.Comment, error)
	DeleteComment(ctx context.Context, id int64) error
}

func NewStorage(logger logger.Logger, config config.Storage, db *dbpg.DB) Storage {
	return postgres.NewStorage(logger, config, db)
}

func ConnectDB(config config.Storage) (*dbpg.DB, error) {

	options := &dbpg.Options{
		MaxOpenConns:    config.MaxOpenConns,
		MaxIdleConns:    config.MaxIdleConns,
		ConnMaxLifetime: config.ConnMaxLifetime,
	}

	db, err := dbpg.New(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.DBName, config.SSLMode), nil, options)
	if err != nil {
		return nil, fmt.Errorf("database driver not found or DSN invalid: %w", err)
	}

	if err := db.Master.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return db, nil

}
