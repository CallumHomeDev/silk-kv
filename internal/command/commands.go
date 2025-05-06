package command

import (
	"strconv"
	"github.com/CallumHomeDev/silk-kv/internal/datastore"
	"github.com/CallumHomeDev/silk-kv/internal/protocol"
)

// handler
var Handlers = map[string]func(store *datastore.Store, args [][]byte) []byte{
	"PING": func(_ *datastore.Store, args [][]byte) []byte { 
		return protocol.FormatSimpleString("PONG") 
	},
	"ECHO": func(_ *datastore.Store, args [][]byte) []byte { 
		return protocol.FormatBulkString(args[1]) 
	},
	"SET": func(store *datastore.Store, args [][]byte) []byte { 
		key := string(args[1])
		value := args[2]
		store.Set(key, value) 
		return protocol.FormatSimpleString("OK")
	},
	"GET": func(store *datastore.Store, args [][]byte) []byte { 
		key := string(args[1])
		if value, oke := store.Get(key); oke { 
			return protocol.FormatBulkString(value) 
		}
		return protocol.FormatNilBulkString()
	},
	// Delete key, return number of keys deleted
	"DEL": func(s *datastore.Store, args [][]byte) []byte { 
		key := string(args[1])
		deleted := s.Delete(key)
		return protocol.FormatInteger(deleted)
	},
	// Check if key exists, return 1 if exists, 0 if not
	"EXISTS": func(s *datastore.Store, args [][]byte) []byte { 
		key := string(args[1])
		if _, ok := s.Get(key); ok {
			return protocol.FormatInteger(1)
		}
		return protocol.FormatInteger(0)
	},
	// Increment the value of a key by 1, if key not exists, set to 1
	"INCR": func(s *datastore.Store, args [][]byte) []byte { 
		key := string(args[1])
		val, _ := s.Get(key)
		n := int64(0)

		if len(val) > 0 {
			x, err := strconv.ParseInt(string(val), 10, 64)
			if err != nil {
				return protocol.FormatError("ERR value is not an integer")
			}
			n = x
		}
		n++
		s.Set(key, []byte(strconv.FormatInt(n, 10)))
		return protocol.FormatInteger(n)
	},
	// Decrement the value of a key by 1, if key not exists, set to -1
	"DECR": func(s *datastore.Store, args [][]byte) []byte { 
		key := string(args[1])
		val, _ := s.Get(key)
		n := int64(0)

		if len(val) > 0 {
			x, err := strconv.ParseInt(string(val), 10, 64)
			if err != nil {
				return protocol.FormatError("ERR value is not an integer")
			}
			n = x
		}
		n--
		s.Set(key, []byte(strconv.FormatInt(n, 10)))
		return protocol.FormatInteger(n)
	},
}