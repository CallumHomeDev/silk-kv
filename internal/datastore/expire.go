// internal/datastore/expire.go
package datastore

import "time"

// Expire sets a TTL in seconds for a key, returns false if key doesn't exist.
func (s *Store) Expire(key string, ttlSeconds int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[key]; !ok {
		return false
	}
	s.expires[key] = time.Now().Add(time.Duration(ttlSeconds) * time.Second)
	return true
}

// TTL returns remaining time-to-live in seconds, or -1 if key missing or no expiration.
func (s *Store) TTL(key string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	exp, ok := s.expires[key]
	if !ok {
		return -1
	}
	remaining := int(exp.Sub(time.Now()).Seconds())
	if remaining < 0 {
		return -1
	}
	return remaining
}