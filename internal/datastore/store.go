// internal/datastore/store.go
package datastore

import (
	"container/list"
	"os"
	"time"
	"sync"
	"github.com/CallumHomeDev/silk-kv/internal/protocol"
)

// Store:keep key-value data and ensure thread-safety
type Store struct {
	mu sync.RWMutex
	data map[string][]byte
	expires map[string]time.Time
	lists map[string][][]byte
	evictionList *list.List
	evictionMap map[string]*list.Element
	capacity int
	aof *os.File
}

// New:create a instance of Store
func New(capacity int, aofFile *os.File) *Store {
	s := &Store{
		data: make(map[string][]byte),
		expires: make(map[string]time.Time),
		lists: make(map[string][][]byte),
		aof: aofFile,
	}
	initEviction(s, capacity)
	return s
}

// Set:save key-value data
func (s *Store) Set(key string, value []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
	delete(s.expires, key) // Remove expiration if set
	s.recordAccess(key)
}

// Get: return value by key
func (s *Store) Get(key string) ([]byte, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// check if the key has expired
	if exp, ok := s.expires[key]; ok {
		if time.Now().After(exp) {
			delete(s.data, key)
			delete(s.expires, key)
			return nil, false
		}
	}

	value, ok := s.data[key]
	if ok {
		s.recordAccess(key)
	}
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

// Exists: check if key exists
func (s *Store) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.data[key]
	return exists
}

// Len: return the number of keys
func (s *Store) Len() int {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return len(s.data)
}

// FlushDB: clear all data
func (s *Store) FlushDB() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.data = make(map[string][]byte)
	s.expires = make(map[string]time.Time)
	s.lists = make(map[string][][]byte)
	s.evictionList.Init()
	s.evictionMap = make(map[string]*list.Element)
}

// AppendCmd write command RESP to AOF file
func (s *Store) AppendCmd(cmd [][]byte) error {
	// convert args to RESP array
	data := protocol.FormatArray(cmd)

	// Write to file
	if _, err := s.aof.Write(data); err != nil {
        return err
    }

	return s.aof.Sync()
}