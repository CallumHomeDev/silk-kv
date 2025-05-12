// datastore/eviction.go

package datastore

import (
	"container/list"
)

// evictionEntry represents an entry in the eviction list.
type evictionEntry struct {
	key string
}

// initEviction initializes the eviction structure.
func initEviction(s *Store, capacity int) {
	s.evictionList = list.New()
	s.evictionMap = make(map[string]*list.Element, capacity)
	s.capacity = capacity
}

// recordAccess update position of the key in the eviction list, handle eviction if necessary.
func (s *Store) recordAccess(key string) {
	if element, ok := s.evictionMap[key]; ok {
		// Move the accessed entry to the front of the list
		s.evictionList.MoveToFront(element)
	} else {
		// Create a new entry
		entry := &evictionEntry{key: key}
		element := s.evictionList.PushFront(entry)
		s.evictionMap[key] = element

		// If the list exceeds the capacity, evict the least recently used entry
		if s.evictionList.Len() > s.capacity {
			// Remove the least recently used entry from the list and map
			tail := s.evictionList.Back()
			old := tail.Value.(*evictionEntry)
			
			// Delete data and expire the entry
			delete(s.data, old.key)
			delete(s.expires, old.key)

			// Delete the entry from the eviction map
			s.evictionList.Remove(tail)
			delete(s.evictionMap, old.key)
		}
	}
}
