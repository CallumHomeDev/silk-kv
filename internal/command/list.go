// internal/command/list.go
package command

import (
    "strconv"
    "github.com/CallumHomeDev/silk-kv/internal/datastore"
    "github.com/CallumHomeDev/silk-kv/internal/protocol"
)

func LPush(store *datastore.Store, args [][]byte) []byte {
    key := string(args[1])
    count := store.LPush(key, args[2:]...)
    return protocol.FormatInteger(int64(count))
}

func RPush(store *datastore.Store, args [][]byte) []byte {
    key := string(args[1])
    count := store.RPush(key, args[2:]...)
    return protocol.FormatInteger(int64(count))
}

func LPop(store *datastore.Store, args [][]byte) []byte {
    key := string(args[1])
    if v, ok := store.LPop(key); ok {
        return protocol.FormatBulkString(v)
    }
    return protocol.FormatNilBulkString()
}

func RPop(store *datastore.Store, args [][]byte) []byte {
    key := string(args[1])
    if v, ok := store.RPop(key); ok {
        return protocol.FormatBulkString(v)
    }
    return protocol.FormatNilBulkString()
}

func LRange(store *datastore.Store, args [][]byte) []byte {
    key := string(args[1])
    start, _ := strconv.Atoi(string(args[2]))
    end, _ := strconv.Atoi(string(args[3]))
    items := store.LRange(key, start, end)
    return protocol.FormatArray(items)
}
