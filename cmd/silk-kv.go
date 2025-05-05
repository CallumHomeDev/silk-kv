package main

import (
	"net"
	"bufio"
	"strings"

	"github.com/CallumHomeDev/silk-kv/internal/protocol"
	"github.com/CallumHomeDev/silk-kv/internal/datastore"
	"github.com/CallumHomeDev/silk-kv/internal/command"
)

func main() {
	store := datastore.New()
	listener, _ := net.Listen("tcp", ":8081")
	for {
		conn, _ := listener.Accept()
		go func (c net.Conn) {
			defer c.Close()
			reader := bufio.NewReader(c)
			for {
				args, err := protocol.ParseRequest(reader)
				if err != nil {
					return
				}

				cmd := strings.ToUpper(string(args[0]))
				handler, exists := command.Handlers[cmd]
				var resp []byte

				if !exists {
					resp = protocol.FormatError("ERR unknown command '" + cmd + "'")
				} else {
					resp = handler(store, args)
				}

				c.Write(resp)
			}
		} (conn)
	}
}