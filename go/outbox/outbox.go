package outbox

import (
	"fmt"

	"github.com/hertzcodes/easy-outbox/go/internal/bindings"
	"github.com/hertzcodes/easy-outbox/go/internal/bindings/pebble"
)

type DBType = string

const (
	DBTypePebble DBType = "pebble"
)

type OutBox struct {
	db bindings.DB
}

func New(db DBType, path string) (*OutBox, error) {
	switch db {
	case DBTypePebble:
		db, err := pebble.New(path)
		if err != nil {
			return nil, err
		}
		return &OutBox{
			db: db,
		}, nil
	default:
		return nil, fmt.Errorf("Invalid database type %s", db)
	}
}

