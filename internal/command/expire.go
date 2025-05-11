// internal/command/expire.go
package command

import (
    "strconv"
    "github.com/CallumHomeDev/silk-kv/internal/datastore"
    "github.com/CallumHomeDev/silk-kv/internal/protocol"
)

func Expire(store *datastore.Store, args [][]byte) []byte {
    key := string(args[1])
    secs, err := strconv.Atoi(string(args[2]))
    if err != nil {
        return protocol.FormatError("ERR value is not an integer")
    }
    if store.Expire(key, secs) {
        return protocol.FormatInteger(1)
    }
    return protocol.FormatInteger(0)
}

func TTL(store *datastore.Store, args [][]byte) []byte {
    key := string(args[1])
    t := store.TTL(key)
    return protocol.FormatInteger(int64(t))
}