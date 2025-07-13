// internal/datastore/list.go
package datastore

// LPush inserts values at the head of the list stored at key, returns new length.
func (s *Store) LPush(key string, values ...[]byte) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	lst := s.lists[key]
	// insert in reverse to preserve order
	for i := len(values) - 1; i >= 0; i-- {
		lst = append([][]byte{values[i]}, lst...)
	}
	s.lists[key] = lst
	return len(lst)
}

// RPush appends values at the tail of the list, returns new length.
func (s *Store) RPush(key string, values ...[]byte) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	lst := s.lists[key]
	lst = append(lst, values...)
	s.lists[key] = lst
	return len(lst)
}

// LPop removes and returns the first element of the list, or nil if empty.
func (s *Store) LPop(key string) ([]byte, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	lst := s.lists[key]
	if len(lst) == 0 {
		return nil, false
	}
	val := lst[0]
	s.lists[key] = lst[1:]
	return val, true
}

// RPop removes and returns the last element of the list, or nil if empty.
func (s *Store) RPop(key string) ([]byte, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	lst := s.lists[key]
	if len(lst) == 0 {
		return nil, false
	}
	idx := len(lst) - 1
	val := lst[idx]
	s.lists[key] = lst[:idx]
	return val, true
}

// LRange returns a slice of elements from start to end (inclusive), supports negative indices.
func (s *Store) LRange(key string, start, end int) [][]byte {
	s.mu.RLock()
	defer s.mu.RUnlock()
	lst := s.lists[key]
	length := len(lst)
	if length == 0 {
		return [][]byte{}
	}
	// handle negative indices
	if start < 0 {
		start = length + start
	}
	if end < 0 {
		end = length + end
	}
	if start < 0 {
		start = 0
	}
	if end >= length {
		end = length - 1
	}
	if start > end || start >= length {
		return [][]byte{}
	}
	return lst[start : end+1]
}