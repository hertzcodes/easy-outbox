package bindings

import "errors"

var (
	ErrNotFound         = errors.New("key not found in database")
	ErrDatabaseNotFound = errors.New("couldn't find database")
)

type DB interface {
	Write(key string, value any) error
	Read(key string) (interface{}, error)
	Delete(key string) error
	ReadBulkKeys(amount int) []string
	PrintMetrics()
}
