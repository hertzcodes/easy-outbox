package outbox

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestPebble(t *testing.T) {

	box, err := New(DBTypePebble, "./testdb_pebble")
	if err != nil {
		panic(err)
	}

	t.Cleanup(func() {
		box.db.PrintMetrics()
		if err := os.RemoveAll("./testdb_pebble"); err != nil {
			t.Log(err)
		}
	})
	t.Run("PopulateKeys", func(t *testing.T) {
		d := "key_%d"
		for i := range 100000 {
			err := box.SetMessage(fmt.Sprintf(d, i), nil)
			if err != nil {
				panic(err)
			}
		}
	})

	t.Run("GetBulkKeys", func(t *testing.T) {
		keys, err := box.GetMessages(100000)
		if err != nil {
			panic(err)
		}
		if len(keys) != 100000 {
			t.Fatalf("expected length of keys to be 100000 got %d", len(keys))
		}
	})

	t.Run("IntervalMode", func(t *testing.T) {
		ch := make(chan []string)
		if err := box.SetInterval(ch, 1*time.Second, 100); err != nil {
			t.Fatal(err)
		}
		var cnt int
		for l := range ch {
			cnt++
			for i := range 100 {
				// t.Log(l[i])
				if err := box.Delete(l[i]); err != nil {
					t.Fatal(err)
				}
			}

			if cnt >= 3 {
				break
			}
		}

		if _, err := box.GetMessages(1000); err == nil {
			t.Fatalf("expected error cause interval is set.")
		}

	})
}
