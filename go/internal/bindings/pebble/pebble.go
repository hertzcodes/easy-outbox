package pebble

import (
	"errors"

	"github.com/cockroachdb/pebble"
	"github.com/hertzcodes/easy-outbox/go/internal/bindings"
	"github.com/hertzcodes/easy-outbox/go/internal/utils"
)

type PebbleStorage struct {
	db *pebble.DB
}

func New(path string) (*PebbleStorage, error) {
	db, err := pebble.Open(path, &pebble.Options{})
	if err != nil {
		return nil, err
	}
	return &PebbleStorage{db: db}, nil
}

func (p *PebbleStorage) Write(key string, value any) error {

	b, err := utils.MarshalGob(value)
	if err != nil {
		return err
	}
	return p.db.Set([]byte(key), b, pebble.Sync)
}

func (p *PebbleStorage) Read(key string) (interface{}, error) {

	value, closer, err := p.db.Get([]byte(key))
	if err != nil {
		if errors.Is(err, pebble.ErrNotFound) {
			return nil, bindings.ErrNotFound
		} else if errors.Is(err, pebble.ErrDBDoesNotExist) {
			return nil, bindings.ErrDatabaseNotFound
		}
		return nil, err
	}
	defer closer.Close()

	var o any
	return o, utils.UnmarshalGob(value, &o)
}

func (p *PebbleStorage) Delete(key string) error {
	return p.db.Delete([]byte(key), pebble.Sync)
}
