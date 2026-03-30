// Package repository defines the data access layer abstraction.
// It provides interfaces and constructors for working with persistent storage,
// isolating the rest of the application from database-specific implementations.
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

// Storage defines the contract for comment persistence.
// It abstracts underlying storage implementation (e.g., Postgres)
// and exposes methods for managing comments and their hierarchy.
type Storage interface {
	CreateComment(ctx context.Context, comment models.Comment) (int64, error)                      // CreateComment stores a new comment and returns its generated ID.
	GetRootComments(ctx context.Context, queryParams models.QueryParams) ([]models.Comment, error) // GetRootComments retrieves top-level comments using pagination and sorting parameters.
	GetCommentTree(ctx context.Context, id int64) ([]models.Comment, error)                        // GetCommentTree returns a full subtree of comments starting from the given root ID.
	DeleteComment(ctx context.Context, id int64) error                                             // DeleteComment removes a comment (and potentially its subtree, depending on implementation).
	Close()                                                                                        // Close releases all underlying resources (e.g., DB connections).
}

// NewStorage constructs a Storage implementation backed by Postgres.
// It wires the logger, configuration, and database connection into the repository layer.
func NewStorage(logger logger.Logger, config config.Storage, db *dbpg.DB) Storage {
	return postgres.NewStorage(logger, config, db)
}

// ConnectDB initializes a Postgres database connection using provided configuration.
// It configures connection pool settings and verifies connectivity via ping.
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
