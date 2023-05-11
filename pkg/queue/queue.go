package queue

import "sync"

type Queue[T comparable] struct {
	sync.RWMutex
	array []T
}

func NewQueue[T comparable](size int) *Queue[T] {
	return &Queue[T]{
		array: make([]T, size),
	}
}

func (s *Queue[T]) GetElement() T {
	s.Lock()
	defer s.Unlock()
	if len(s.array) != 0 {
		return s.array[0]
	}
	var result T
	return result
}

func (s *Queue[T]) Push(value T) {
	s.RLock()
	defer s.RUnlock()
	s.array = append(s.array, value)
}

func (s *Queue[T]) Pop() T {
	s.RLock()
	defer s.RUnlock()
	var result T
	result, s.array = s.array[0], s.array[1:]
	return result
}
