package cache

import (
	"container/list"
	"sync"
	"time"
)

type Store struct {
	mu       *sync.Mutex
	store    map[string]*list.Element
	ll       *list.List
	maxItems int
}

func NewStore(maxItems int) *Store {
	return &Store{
		mu:       &sync.Mutex{},
		store:    make(map[string]*list.Element),
		ll:       &list.List{},
		maxItems: maxItems,
	}
}

type Node struct {
	key    string
	value  any
	expire time.Time
}