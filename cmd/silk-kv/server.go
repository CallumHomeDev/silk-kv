package main

import (
	"fmt"
	"log"

	"github.com/panjf2000/gnet"
    "github.com/CallumHomeDev/silk-kv/internal/datastore"
    "github.com/CallumHomeDev/silk-kv/internal/command"
    "github.com/CallumHomeDev/silk-kv/internal/protocol"
)

type SilkKVServer struct {
	*gnet.EventServer
	store *datastore.Store
}

func newServer(capacity int) *SilkKVServer {
	return &SilkKVServer{
		store: datastore.NewStore(capacity),
	}
} 

// OnInitComplete is called when the server is initialized.
func (s *SilkKVServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("🚀 gnet server started on %s (multicore=%t)", srv.Addr.String(), srv.Multicore)
	return
}

// React is called when a new connection is established.
func (s *SilkKVServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	// frame is the raw data received from the client
	args, err := protocol.ParseRequest(frame)
	if err != nil {
		out = protocol.FormatError(err.Error())
		return out, gnet.Close
	}

	cmd := strings.ToUpper(string(args[0]))
	handler, ok := command.Commands[cmd]
	if !ok {
		out = protocol.FormatError(fmt.Sprintf("unknown command '%s'", cmd))
	} else {
		out = handler(s.store, args)
	}
	return out, gnet.None
}

func main() {
	addr := "tcp://:8081"
	server := newServer(10000) // Set the capacity to 10000

	// Run 4 reactor threads, enable multicore mode
	log.Fatal(gnet.Serve(server, addr, gnet.WithMulticore(true), gnet.WithNumEventLoop(4)))
}