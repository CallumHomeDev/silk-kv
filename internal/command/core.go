// internal/command/core.go
package command

import (
    "strconv"
    "github.com/CallumHomeDev/silk-kv/internal/datastore"
    "github.com/CallumHomeDev/silk-kv/internal/protocol"
)

// PING handler
func Ping(_ *datastore.Store, args [][]byte) []byte {
    return protocol.FormatSimpleString("PONG")
}

// ECHO handler
func Echo(_ *datastore.Store, args [][]byte) []byte {
    return protocol.FormatBulkString(args[1])
}

// Set, Get, Del, Exists, Incr, Decr handlers
func Set(store *datastore.Store, args [][]byte) []byte {
    key := string(args[1])
    store.Set(key, args[2])
    return protocol.FormatSimpleString("OK")
}

func Get(store *datastore.Store, args [][]byte) []byte {
    key := string(args[1])
    if v, ok := store.Get(key); ok {
        return protocol.FormatBulkString(v)
    }
    return protocol.FormatNilBulkString()
}

func Del(store *datastore.Store, args [][]byte) []byte {
    key := string(args[1])
    n := store.Delete(key)
    return protocol.FormatInteger(int64(n))
}

func Exists(store *datastore.Store, args [][]byte) []byte {
    key := string(args[1])
    if _, ok := store.Get(key); ok {
        return protocol.FormatInteger(1)
    }
    return protocol.FormatInteger(0)
}

func Incr(store *datastore.Store, args [][]byte) []byte {
    key := string(args[1])
    v, _ := store.Get(key)
    n := int64(0)
    if len(v) > 0 {
        x, err := strconv.ParseInt(string(v), 10, 64)
        if err != nil {
            return protocol.FormatError("ERR value is not an integer")
        }
        n = x
    }
    n++
    store.Set(key, []byte(strconv.FormatInt(n, 10)))
    return protocol.FormatInteger(n)
}

func Decr(store *datastore.Store, args [][]byte) []byte {
    key := string(args[1])
    v, _ := store.Get(key)
    n := int64(0)
    if len(v) > 0 {
        x, err := strconv.ParseInt(string(v), 10, 64)
        if err != nil {
            return protocol.FormatError("ERR value is not an integer")
        }
        n = x
    }
    n--
    store.Set(key, []byte(strconv.FormatInt(n, 10)))
    return protocol.FormatInteger(n)
}