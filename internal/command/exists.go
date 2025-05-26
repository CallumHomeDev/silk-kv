// internal/command/exists.go

package command

import (
	"parser"
	"github.com/CallumHomeDev/silk-kv/internal/datastore"
)

func Exists(store *datastore.Store, parts []parser.Token) parser.Response {
	if len(parts) != 2 {
		return parser.Error("ERR wrong number of arguments for 'exists' command")
	}

	key := parts[1].String()
	if store.Exists(key) {
		return parser.NewInteger(1)
	}

	return parser.NewInteger(0)
}