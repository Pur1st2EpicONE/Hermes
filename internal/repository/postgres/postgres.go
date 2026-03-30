// Package postgres provides a Postgres-backed implementation
// of the repository.Storage interface. It contains SQL queries,
// retry logic, and low-level data access operations.
package postgres

import (
	"Hermes/internal/config"
	"Hermes/internal/logger"

	"github.com/wb-go/wbf/dbpg"
)

// Storage implements the repository.Storage interface using Postgres.
// It encapsulates database access, logging, and configuration for retries.
type Storage struct {
	db     *dbpg.DB       // db is the underlying database connection pool.
	logger logger.Logger  // logger is used for logging database-related events.
	config config.Storage // config holds DB and retry configuration.
}

// NewStorage constructs a new Postgres-backed Storage instance.
func NewStorage(logger logger.Logger, config config.Storage, db *dbpg.DB) *Storage {
	return &Storage{db: db, logger: logger, config: config}
}

// Close gracefully shuts down the database connection.
// It logs the result of the shutdown operation.
func (s *Storage) Close() {
	if err := s.db.Master.Close(); err != nil {
		s.logger.LogError("postgres — failed to close properly", err, "layer", "repository.postgres")
	} else {
		s.logger.LogInfo("postgres — database closed", "layer", "repository.postgres")
	}
}

// DB exposes the underlying dbpg.DB instance.
// Intended for internal use where direct DB access is required.
func (s *Storage) DB() *dbpg.DB {
	return s.db
}
