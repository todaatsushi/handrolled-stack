package cache_test

import (
	"errors"
	"testing"
	"time"

	"github.com/todaatsushi/handrolled-cache/internal/cache"
)

type c struct{}

func (clock c) Now() time.Time {
	t, err := time.Parse(time.RFC3339, "2069-04-20T14:00:00Z")
	if err != nil {
		panic(err)
	}
	return t
}

func (clock c) CalcExpires(ttl int) time.Time {
	t, err := time.Parse(time.RFC3339, "2069-04-20T14:59:50Z")
	if err != nil {
		panic(err)
	}
	return t.Add(time.Second * time.Duration(ttl))
}

func TestSet(t *testing.T) {
	t.Run("Set value", func(t *testing.T) {
		s := cache.NewStore(1, c{})

		if s.NumItems != 0 {
			t.Errorf("Expected %d items, got %d", 0, s.NumItems)
		}

		_, err := s.Set("key", 420, 10)
		if err != nil {
			t.Error(err)
		}

		if s.NumItems != 1 {
			t.Errorf("Expected %d items, got %d", 1, s.NumItems)
		}
	})

	t.Run("Update existing value", func(t *testing.T) {
		s := cache.NewStore(1, c{})
		expires, err := s.Set("key", 420, 0)
		if err != nil {
			t.Fatal(err)
		}

		newExpires, err := s.Set("key", 420, 1000)
		if err != nil {
			t.Fatal(err)
		}

		if newExpires.Compare(expires) != 1 {
			t.Fatal("Expected new expires to be after old expires.")
		}
	})
}
