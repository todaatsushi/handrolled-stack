package cache

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

type Clock interface {
	Now() time.Time
	Expired(t time.Time) bool
}

type Store struct {
	mu       *sync.Mutex
	store    map[string]*list.Element
	ll       *list.List
	maxItems uint64 // 0 == unlimited
	NumItems uint64
	C        Clock
}

func (s *Store) Set(key string, value string, expires time.Time) (exp time.Time, err error) {
	if expires.Compare(s.C.Now()) == -1 {
		return expires, errors.New("Expiry can't be in the past.")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.store[key]
	if ok {
		s.ll.MoveToFront(item)
		node := item.Value.(*Node)
		node.Expire = expires
		return node.Expire, nil
	}

	node := &Node{
		key, []byte(value), expires,
	}
	s.ll.PushFront(node)
	s.NumItems++

	s.store[key] = s.ll.Front()

	if s.NumItems > s.maxItems && s.maxItems != 0 {
		last := s.ll.Back()

		delete(s.store, last.Value.(*Node).Key)
		s.ll.Remove(last)
		s.NumItems--
	}
	return node.Expire, nil
}

func (s *Store) Get(key string) (value []byte, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.store[key]
	if !ok {
		return nil, errors.New("Value doesn't exist.")
	}

	node := item.Value.(*Node)
	if s.C.Expired(node.Expire) {
		return nil, errors.New("Expired.")
	}

	s.ll.MoveToFront(item)
	return node.Value, nil
}

func NewStore(maxItems uint64, c Clock) *Store {
	return &Store{
		mu:       &sync.Mutex{},
		store:    make(map[string]*list.Element),
		ll:       list.New(),
		maxItems: maxItems,
		NumItems: 0,
		C:        c,
	}
}

type Node struct {
	Key    string
	Value  []byte
	Expire time.Time
}
