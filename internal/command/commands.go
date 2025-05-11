// commands/commands.go

package command

import (
    "github.com/CallumHomeDev/silk-kv/internal/datastore"
)

// Handlers maps Redis command names to their handler functions.
// This file aggregates all individual handlers defined in other files.
var Handlers = map[string]func(store *datastore.Store, args [][]byte) []byte{
    // Core commands
    "PING":     Ping,
    "ECHO":     Echo,
    "SET":      Set,
    "GET":      Get,
    "DEL":      Del,
    "EXISTS":   Exists,
    "INCR":     Incr,
    "DECR":     Decr,

    // Expiration commands
    "EXPIRE":   Expire,
    "TTL":      TTL,

    // List commands
    "LPUSH":    LPush,
    "RPUSH":    RPush,
    "LPOP":     LPop,
    "RPOP":     RPop,
    "LRANGE":   LRange,
}