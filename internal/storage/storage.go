// Package storage represents abstract storage.
package storage

import (
	"database/sql"

	"github.com/evgenytr/metrics.git/internal/interfaces"
	"github.com/evgenytr/metrics.git/internal/storage/database"
	"github.com/evgenytr/metrics.git/internal/storage/memstorage"
)

// NewStorage returns Storage interface to memstorage or database depending on params.
func NewStorage(db *sql.DB, fileStoragePath string) interfaces.Storage {
	if db != nil {
		return database.NewStorage(db)
	}
	return memstorage.NewStorage(fileStoragePath)
}
