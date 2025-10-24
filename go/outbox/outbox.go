package outbox

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/hertzcodes/easy-outbox/go/internal/bindings"
	"github.com/hertzcodes/easy-outbox/go/internal/bindings/pebble"
)

type DBType = string

const (
	DBTypePebble DBType = "pebble"
)

type OutBox struct {
	db             bindings.DB
	r_ch           chan []string // the channel for reads
	d_ch           chan string   // the channel for deletions
	interval_count int
	interval       time.Duration
	once           sync.Once
}

func New(db DBType, path string) (*OutBox, error) {
	switch db {
	case DBTypePebble:
		db, err := pebble.New(path)
		if err != nil {
			return nil, err
		}
		return &OutBox{
			db:   db,
			d_ch: make(chan string),
		}, nil
	default:
		return nil, fmt.Errorf("Invalid database type %s", db)
	}
}

func (o *OutBox) SetInterval(ch chan []string, t time.Duration, cnt int) error {
	o.once.Do(func() {
		if ch == nil {
			panic("nil channel was given")
		}
		o.r_ch = ch
		o.interval = t
		o.interval_count = cnt
		go func() {
			for {
				o.r_ch <- o.db.ReadBulkKeys(cnt)
				time.Sleep(o.interval)
			}
		}()
	})
	return nil
}

func (o *OutBox) GetMessages(count int) ([]string, error) {
	if o.r_ch != nil {
		return nil, errors.New("Can't use Bulk Message fetch when interval is on")
	}
	return o.db.ReadBulkKeys(count), nil
}

func (o *OutBox) SetMessage(key string, value interface{}) error {
	return o.db.Write(key, value)
}

// TODO: this has to be done through d_ch for bulk deletions
func (o *OutBox) Delete(key string) error {
	return o.db.Delete(key)
}
