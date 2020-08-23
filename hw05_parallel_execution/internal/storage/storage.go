package storage

type Task func() error

type Interface interface {
	Next() Task
}

type Storage struct {
	buffer []Task
	index  int
}

// New create and return new instance of element storage.
func New(dataCopy []Task) Interface {
	return &Storage{
		buffer: dataCopy,
	}
}

// Next returns next element from buffer and move forvard internal index,
// returns nil if there is no next element.
func (s *Storage) Next() Task {
	var result Task
	if s.index < len(s.buffer) {
		result = s.buffer[s.index]
		s.index++
	}

	return result
}
