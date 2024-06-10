package service

import "sync"

type ThreadSafeMap[V any] struct {
	mutex *sync.Mutex
	data  map[string]V
}

func (m *ThreadSafeMap[V]) Set(key string, value V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data[key] = value
}

func (m *ThreadSafeMap[V]) Get(key string) V {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.data[key]
}

func (m *ThreadSafeMap[V]) Modify(fx func(map[string]V) map[string]V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data = fx(m.data)
}

func (m *ThreadSafeMap[V]) Unset(key string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.data, key)
}

func (m *ThreadSafeMap[V]) Contain(key string) bool {
	{
		m.mutex.Lock()
		defer m.mutex.Unlock()
		_, ok := m.data[key]
		return ok
	}
}

func NewThreadSafeMap[V any]() *ThreadSafeMap[V] {
	return &ThreadSafeMap[V]{
		mutex: &sync.Mutex{},
		data:  make(map[string]V),
	}
}
