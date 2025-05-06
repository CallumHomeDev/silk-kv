package datastore

import "sync"

// Store:keep key-value data and ensure thread-safety
type Store struct {
	mu sync.RWMutex
	data map[string][]byte
}

// New:create a instance of Store
func New() *Store {
	return &Store{
		data: make(map[string][]byte),
	}
}

// Set:save key-value data
func (s *Store) Set(key string, value []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// Get: return value by key
func (s *Store) Get(key string) ([]byte, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.data[key]
	return value, ok
}

// Delete: delete key-value data
func (s *Store) Delete(key string) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[key]; ok {
		delete(s.data, key)
		return 1
	}
	return 0
}