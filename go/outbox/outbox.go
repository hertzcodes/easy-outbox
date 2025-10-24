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
	r_ch           chan string // the channel for reads
	d_ch           chan string // the channel for deletions
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

// WithStream sets the channel for streaming keys. Must set a buffered channel to avoid blocking and too much writes.
//
// Example:
//
//     // Create a buffered channel
//     ch := make(chan string, 100)
//     
//     // Set up the stream
//     if err := outbox.WithStream(ch); err != nil {
//         t.Fatal(err)
//     }
//     
//     // Process keys
//     for key := range ch {
//         fmt.Println(key)
//     }
func (o *OutBox) WithStream(ch chan string) error {
	o.once.Do(func() {
		if ch == nil {
			panic("nil channel was given")
		}
		if cap(ch) == 0 {
			panic("channel is not buffered")
		}
		o.r_ch = ch
		o.db.StreamKeys(o.r_ch)
	})
	return nil
}

func (o *OutBox) GetMessages(count int) ([]string, error) {
	if o.r_ch != nil {
		return nil, errors.New("Can't use Bulk Message fetch when stream is on")
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
