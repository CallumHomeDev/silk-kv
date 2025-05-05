package command

import (
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
}