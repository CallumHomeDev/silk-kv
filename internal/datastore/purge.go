// internal/datastore/purge.go

package datastore

import (
	"time"
)

// Purge clears all data from the store, including keys and their expirations.
func (s *Store) PurgeExpiredKeys() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for key, exp := range s.expires {
		if exp.Before(now) {
			delete(s.data, key)
			delete(s.expires, key)
		}
	}
}