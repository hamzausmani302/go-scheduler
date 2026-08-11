//go:build !cgo
// +build !cgo

package storage

import "errors"

var errCgoRequired = errors.New("sqlite3 storage requires building with CGO_ENABLED=1 (and a C compiler)")

// Sqlite3Config is the config structure holding information about sqlite db.
type Sqlite3Config struct {
	DbName string
}

// Sqlite3Storage is a stub used when the binary is built without cgo.
type Sqlite3Storage struct{}

// NewSqlite3Storage returns a new instance of Sqlite3Storage.
func NewSqlite3Storage(config Sqlite3Config) Sqlite3Storage {
	return Sqlite3Storage{}
}

// Connect always fails; sqlite3 storage is unavailable without cgo.
func (sqlite *Sqlite3Storage) Connect() error {
	return errCgoRequired
}

// Close is a no-op.
func (sqlite Sqlite3Storage) Close() error {
	return nil
}

// Initialize always fails; sqlite3 storage is unavailable without cgo.
func (sqlite *Sqlite3Storage) Initialize() error {
	return errCgoRequired
}

// Add always fails; sqlite3 storage is unavailable without cgo.
func (sqlite Sqlite3Storage) Add(task TaskAttributes) error {
	return errCgoRequired
}

// Remove always fails; sqlite3 storage is unavailable without cgo.
func (sqlite Sqlite3Storage) Remove(task TaskAttributes) error {
	return errCgoRequired
}

// Fetch always fails; sqlite3 storage is unavailable without cgo.
func (sqlite Sqlite3Storage) Fetch() ([]TaskAttributes, error) {
	return nil, errCgoRequired
}
