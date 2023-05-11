package cache

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"urinal/pkg/queue"
)

type structValueAndTime[V any] struct {
	Value V
	Time  int64
}

type Cache[T comparable, V any] struct {
	sync.Mutex
	cache   map[T]structValueAndTime[V]
	queue   *queue.Queue[T]
	timeout time.Duration
}

func NewCache[T comparable, V any](timeout time.Duration) *Cache[T, V] {
	queue := queue.NewQueue[T](100)
	return &Cache[T, V]{
		cache:   make(map[T]structValueAndTime[V], 100),
		queue:   queue,
		timeout: timeout,
	}
}

func (c *Cache[T, V]) CheckTime() error {
	element, ok := c.cache[c.queue.GetElement()]
	if !ok {
		return errors.New("Not found element")
	}
	if time.Now().Local().UnixMilli()-element.Time > c.timeout.Milliseconds() {
		fmt.Println("DELETE")
		delete(c.cache, c.queue.Pop())
		return nil
	}
	return nil
}

func (c *Cache[T, V]) AddElement(key T, value V) {
	c.Lock()
	_, ok := c.cache[key]
	if ok {
		c.queue.Push(c.queue.Pop())
		st := c.cache[key]
		st.Time = time.Now().UnixMilli() + c.timeout.Milliseconds()
		st.Value = value
		c.Unlock()
		return
	}
	c.queue.Push(key)
	c.cache[key] = structValueAndTime[V]{
		Value: value,
		Time:  time.Now().UnixMilli() + c.timeout.Milliseconds(),
	}
	c.Unlock()
}

func (c *Cache[T, V]) Get(key T) (V, bool) {
	c.Lock()
	defer c.Unlock()
	value, ok := c.cache[key]
	if ok {
		c.queue.Push(c.queue.Pop())
		value.Time = time.Now().UnixMilli() + c.timeout.Milliseconds()
		c.cache[key] = value
	}
	return value.Value, ok
}
