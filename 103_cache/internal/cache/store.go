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
	maxItems uint64 // 0 == unlimited
	NumItems uint64
}

func NewStore(maxItems uint64) *Store {
	return &Store{
		mu:       &sync.Mutex{},
		store:    make(map[string]*list.Element),
		ll:       list.New(),
		maxItems: maxItems,
		NumItems: 0,
	}
}

type Node struct {
	Key    string
	Value  any
	Expire time.Time
}
