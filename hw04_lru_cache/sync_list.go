package hw04_lru_cache //nolint:golint,stylecheck

import (
	"sync"
)

type syncList struct {
	l  List
	mu *sync.RWMutex
}

func NewSyncList(l List) List {
	return &syncList{
		l:  l,
		mu: &sync.RWMutex{},
	}
}

func (sl syncList) Len() int {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	return sl.l.Len()
}

func (sl syncList) Front() *Item {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	return sl.l.Front()
}

func (sl syncList) Back() *Item {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	return sl.l.Back()
}

func (sl *syncList) PushFront(v interface{}) *Item {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	return sl.l.PushFront(v)
}

func (sl *syncList) PushBack(v interface{}) *Item {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	return sl.l.PushBack(v)
}

func (sl *syncList) Remove(i *Item) {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	sl.l.Remove(i)
}

func (sl *syncList) MoveToFront(i *Item) {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	sl.l.MoveToFront(i)
}
