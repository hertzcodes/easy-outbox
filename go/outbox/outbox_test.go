package outbox

import (
	"fmt"
	"os"
	"testing"
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
		for i := range 10 {
			err := box.SetMessage(fmt.Sprintf(d, i), nil)
			if err != nil {
				panic(err)
			}
		}
	})

	t.Run("GetBulkKeys", func(t *testing.T) {
		keys, err := box.GetMessages(10)
		if err != nil {
			panic(err)
		}
		if len(keys) != 10 {
			t.Fatalf("expected length of keys to be 10 got %d", len(keys))
		}
	})

	// ignore this test for now
	// t.Skip()
	t.Run("StreamMode", func(t *testing.T) {
		ch := make(chan string, 10)
		if err := box.WithStream(ch); err != nil {
			t.Fatal(err)
		}
		var cnt int
		// blocks until a new key is written to the channel
		for key := range ch {
			cnt++
			box.Delete(key)
			t.Log(key)
			if cnt >= 10 {
				break
			}
		}

		if _, err := box.GetMessages(100); err == nil {
			t.Fatalf("expected error cause interval is set.")
		}

	})
}
