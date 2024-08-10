package cache

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

type Clock interface {
	Now() time.Time
	CalcExpires(ttl int) time.Time
}

type Store struct {
	mu       *sync.Mutex
	store    map[string]*list.Element
	ll       *list.List
	maxItems uint64 // 0 == unlimited
	NumItems uint64
	c        Clock
}

func (s *Store) Set(key string, value any, ttl int) (expires time.Time, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e := s.c.CalcExpires(ttl)
	item, ok := s.store[key]
	if ok {
		s.ll.MoveToFront(item)
		node := item.Value.(*Node)
		node.Expire = e
		return node.Expire, nil
	}

	node := &Node{
		key, value, e,
	}
	s.ll.PushFront(node)
	s.NumItems++

	s.store[key] = s.ll.Front()

	if s.NumItems > s.maxItems && s.maxItems != 0 {
		last := s.ll.Back()
		s.ll.Remove(last)
		s.NumItems--
	}
	return node.Expire, nil
}

func (s *Store) Get(key string) (value any, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.store[key]
	if !ok {
		return nil, errors.New("Value doesn't exist.")
	}

	node := item.Value.(*Node)
	now := s.c.Now()

	if node.Expire.Unix() < now.Unix() {
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
		c:        c,
	}
}

type Node struct {
	Key    string
	Value  any
	Expire time.Time
}
