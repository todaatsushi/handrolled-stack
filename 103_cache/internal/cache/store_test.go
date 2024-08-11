package cache_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/todaatsushi/handrolled-cache/internal/cache"
)

type c struct {
	expired bool
}

func (clock c) Now() time.Time {
	t, err := time.Parse(time.RFC3339, "2069-04-20T15:00:00Z")
	if err != nil {
		panic(err)
	}
	return t
}

func (clock c) Before() time.Time {
	t, err := time.Parse(time.RFC3339, "2069-04-20T14:00:00Z")
	if err != nil {
		panic(err)
	}
	return t
}

func (clock c) Future() time.Time {
	t, err := time.Parse(time.RFC3339, "2069-04-20T16:00:00Z")
	if err != nil {
		panic(err)
	}
	return t
}

func (clock c) Expired(t time.Time) bool {
	return clock.expired
}

var clock = c{false}

func TestSet(t *testing.T) {
	t.Run("Set value", func(t *testing.T) {
		s := cache.NewStore(1, clock)

		if s.NumItems != 0 {
			t.Errorf("Expected %d items, got %d", 0, s.NumItems)
		}

		_, err := s.Set("key", "420", clock.Now())
		if err != nil {
			t.Error(err)
		}

		if s.NumItems != 1 {
			t.Errorf("Expected %d items, got %d", 1, s.NumItems)
		}
	})

	t.Run("Past expiry", func(t *testing.T) {
		s := cache.NewStore(1, clock)
		_, err := s.Set("key", "420", clock.Before())
		if err == nil {
			t.Fatal("Expected err got nil.")
		}

		expected := errors.New("Expiry can't be in the past.").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("Update existing value", func(t *testing.T) {
		s := cache.NewStore(1, clock)
		expires, err := s.Set("key", "420", clock.Now())
		if err != nil {
			t.Fatal(err)
		}

		newExpires, err := s.Set("key", "420", clock.Future())
		if err != nil {
			t.Fatal(err)
		}

		if newExpires.Compare(expires) != 1 {
			t.Fatal("Expected new expires to be after old expires.")
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("Get stored value", func(t *testing.T) {
		expected := "420"
		s := cache.NewStore(1, clock)
		_, err := s.Set("key", expected, clock.Now())
		if err != nil {
			t.Fatal(err)
		}

		data, err := s.Get("key")
		if err != nil {
			t.Fatal(err)
		}

		if data == nil {
			t.Fatal("Expected value, got nil")
		}

		actual := string(data)
		if actual != expected {
			t.Errorf("Expected %s, got %s", expected, actual)
		}
	})

	t.Run("Get non stored value", func(t *testing.T) {
		s := cache.NewStore(1, clock)

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
		expiredClock := c{true}

		s := cache.NewStore(1, expiredClock)

		_, err := s.Set("key", "420", clock.Now())
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

func TestCache(t *testing.T) {
	t.Run("Test eviction", func(t *testing.T) {
		var err error

		// Create store with N values
		s := cache.NewStore(2, c{})

		// Store N values with specific TTL / expires X
		for i := range 2 {
			_, err = s.Set(fmt.Sprint(i), fmt.Sprint(i), clock.Now())
			if err != nil {
				t.Fatal(err)
			}
		}

		// Store next value
		_, err = s.Set("2", "2", clock.Now())
		if err != nil {
			t.Fatal(err)
		}

		// Refetch first value - expires should not be X as it was evicted.
		_, err = s.Get("0")
		if err == nil {
			t.Fatal("Expected err, got nil")
		}

		expected := errors.New("Value doesn't exist.").Error()
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})
}
