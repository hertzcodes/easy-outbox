package pebble

import (
	"errors"
	"fmt"

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

	b, err := utils.Marshal(value)
	if err != nil {
		return err
	}
	return p.db.Set([]byte(key), b, pebble.NoSync)
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

	// TODO: move storage to generic ones and fix this
	var o any
	return o, utils.Unmarshal(value, &o)
}

func (p *PebbleStorage) ReadBulkKeys(amount int) []string {
	iter, err := p.db.NewIter(&pebble.IterOptions{
		LowerBound: nil,
	})
	if err != nil {
		// log the error
		return nil
	}
	defer iter.Close()

	var keys []string
	cnt := 0
	for iter.First(); iter.Valid() && cnt < amount; iter.Next() {
		keys = append(keys, string(iter.Key()))
		cnt++
	}
	return keys
}

func (p *PebbleStorage) Delete(key string) error {
	return p.db.Delete([]byte(key), pebble.NoSync)
}

func (p *PebbleStorage) PrintMetrics() {
	fmt.Println(p.db.Metrics())
}
