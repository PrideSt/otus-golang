package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int                      // длина списка
	Front() *Item                  // первый Item
	Back() *Item                   // последний Item
	PushFront(v interface{}) *Item // добавить значение в начало
	PushBack(v interface{}) *Item  // добавить значение в конец
	Remove(i *Item)                // удалить элемент
	MoveToFront(i *Item)           // переместить элемент в начало
}

type Item struct {
	Next  *Item
	Prev  *Item
	Value interface{}
}

type list struct {
	first *Item
	last  *Item
	len   int
}

// NewList creates new instance of list.
func NewList() List {
	return &list{
		first: nil,
		last:  nil,
		len:   0,
	}
}

// Len returns count of elements in list.
func (l list) Len() int {
	return l.len
}

// Front returns the first element from list.
func (l list) Front() *Item {
	return l.first
}

// Back returns the last element from list.
func (l list) Back() *Item {
	return l.last
}

// PushFront add new element at the begin of list.
func (l *list) PushFront(v interface{}) *Item {
	i := &Item{
		Next:  l.first,
		Prev:  nil,
		Value: v,
	}

	if l.first == nil {
		l.last = i
	} else {
		l.first.Prev = i
	}

	l.first = i
	l.len++

	return i
}

// PushBack add new element at the end of list.
func (l *list) PushBack(v interface{}) *Item {
	i := &Item{
		Next:  nil,
		Prev:  l.last,
		Value: v,
	}

	if l.last == nil {
		l.first = i
	} else {
		l.last.Next = i
	}

	l.last = i
	l.len++

	return i
}

// Remove delete element i from list.
func (l *list) Remove(i *Item) {
	if l.len < 2 {
		l.first = nil
		l.last = nil
		l.len = 0

		return
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.first = i.Next
		i.Next.Prev = nil
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.last = i.Prev
		i.Prev.Next = nil
	}

	l.len--
}

// MoveToFront move element i from current position of list to begin.
func (l *list) MoveToFront(i *Item) {
	if l.len < 2 {
		return
	}

	// i is the first element, nothing to do
	if i.Prev == nil {
		return
	}

	// remove item, bind left and right elements together
	i.Prev.Next = i.Next
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		// remove last element
		l.last = i.Prev
	}

	// place element i at begin
	// we can use PushFront, but don't want to copy value and create new Item node
	i.Prev = nil
	i.Next = l.first
	l.first.Prev = i
	l.first = i
}
