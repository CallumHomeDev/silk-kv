	package main

	import (
		"fmt"
		"log"
		"os"
		"strings"
		"time"
		"io"
	    "strconv"

		"github.com/panjf2000/gnet"
		"github.com/CallumHomeDev/silk-kv/internal/datastore"
		"github.com/CallumHomeDev/silk-kv/internal/command"
		"github.com/CallumHomeDev/silk-kv/internal/protocol"
	)

	type SilkKVServer struct {
		*gnet.EventServer
		store *datastore.Store
	}

	func newServer(capacity int, aofFile *os.File) *SilkKVServer {
		return &SilkKVServer{
			store: datastore.New(capacity, aofFile),
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
		args, _, err := protocol.ParseFrame(frame)
		if err != nil {
			out = protocol.FormatError(err.Error())
			return out, gnet.Close
		}

		cmd := strings.ToUpper(string(args[0]))
		handler, ok := command.Handlers[cmd]
		if !ok {
			out = protocol.FormatError(fmt.Sprintf("unknown command '%s'", cmd))
		} else {
			out = handler(s.store, args)
			if err := s.store.AppendCmd(args); err != nil {
				log.Printf("AOF write error: %v", err)
			}
		}
		return out, gnet.None
	}

	func main() {
		aofFile, err := os.OpenFile("appendonly.aof",
			os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			log.Fatalf("failed to open AOF file: %v", err)
		}
		defer aofFile.Close()
		
		addr := "tcp://:8080"
		server := newServer(10000, aofFile) // Set the capacity to 10000

		// Replay AOF to store
		data, err := os.ReadFile("appendonly.aof")
		if err != nil {
			log.Fatal("failed to read AOF file: %v", err)
		}

		offset := 0
		for offset < len(data) {
			args, n, err := protocol.ParseFrame(data[offset:])
			if err != nil {
				log.Printf("AOF parse error at offset %d: %v", offset, err)
				break
			}

			switch cmd := strings.ToUpper(string(args[0])); cmd {
			case "SET":
				server.store.Set(string(args[1]), args[2])
			case "DEL":
				server.store.Delete(string(args[1]))
			case "EXPIRE":
				sec, _ := strconv.Atoi(string(args[2]))
				server.store.Expire(string(args[1]), sec)
			}
			offset += n
		}

		if _, err := aofFile.Seek(0, io.SeekEnd); err != nil {
			log.Fatalf("failed to seek AOF file: %v", err)
		}	

		// goroutine to periodically clean up expired keys
		go func() {
			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			for range ticker.C {
				server.store.PurgeExpiredKeys()
			}
		} ()

		// Run 4 reactor threads, enable multicore mode
		log.Fatal(gnet.Serve(server, addr, gnet.WithMulticore(true), gnet.WithNumEventLoop(4)))
	}