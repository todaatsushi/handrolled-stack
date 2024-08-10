package cache_test

import (
	"errors"
	"testing"
	"time"

	"github.com/todaatsushi/handrolled-cache/internal/cache"
)

type c struct{}

func (clock c) Now() time.Time {
	t, err := time.Parse(time.RFC3339, "2069-04-20T15:00:00Z")
	if err != nil {
		panic(err)
	}
	return t
}

func (clock c) CalcExpires(ttl int) time.Time {
	t, err := time.Parse(time.RFC3339, "2069-04-20T14:50:00Z")
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

	t.Run("Negative ttl", func(t *testing.T) {
		s := cache.NewStore(1, c{})
		_, err := s.Set("key", 420, -1)
		if err == nil {
			t.Fatal("Expecting err, got nil")
		}

		expected := errors.New("TTL can't be negative.").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("Get stored value", func(t *testing.T) {
		expected := 420
		s := cache.NewStore(1, c{})
		_, err := s.Set("key", expected, 3600)
		if err != nil {
			t.Fatal(err)
		}

		actual, err := s.Get("key")
		if err != nil {
			t.Fatal(err)
		}

		if actual == nil {
			t.Fatal("Expected value, got nil")
		}

		if actual != expected {
			t.Errorf("Expected %d, got %d", expected, actual)
		}
	})

	t.Run("Get non stored value", func(t *testing.T) {
		s := cache.NewStore(0, c{})

		v, err := s.Get("nonexistent")

		if err == nil {
			t.Fatal("Expected err, got nill")
		}

		if v != nil {
			t.Fatalf("Expected nil, got '%s'", v)
		}

		expected := errors.New("Value doesn't exist.").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("Get expired value", func(t *testing.T) {
		s := cache.NewStore(1, c{})
		_, err := s.Set("key", 420, 0)
		if err != nil {
			t.Fatal(err)
		}

		_, err = s.Get("key")
		if err == nil {
			t.Fatal("Expecting err, got nil.")
		}

		expected := errors.New("Expired.").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})
}
