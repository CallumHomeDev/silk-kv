package main

import (
    "bufio"
    "context"
    "fmt"
    "log"
    "net"
    "os"
	"os/signal"
    "strings"
    "sync"
    "syscall"

    "github.com/CallumHomeDev/silk-kv/internal/command"
    "github.com/CallumHomeDev/silk-kv/internal/datastore"
    "github.com/CallumHomeDev/silk-kv/internal/protocol"
)

func main() {
    // Initialize in-memory data store
    store := datastore.New()

    // Start TCP server
    listener, err := net.Listen("tcp", ":8081")
    if err != nil {
        log.Fatalf("Error starting server: %v", err)
    }
    defer listener.Close()
    log.Println("▶︎ Server listening on :8081")

    // Context and Waitgroup for graceful shutdown
    ctx, cancel := context.WithCancel(context.Background())
    var wg sync.WaitGroup

    // Signal handling (SIGINT, SIGTERM)
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigCh
        log.Println("Shutdown signal received, shutting down...")
        cancel()
        listener.Close()
    }()

    // Accept loop
    for {
        conn, err := listener.Accept()
        if err != nil {
            select {
            case <-ctx.Done():
                // Graceful shutdown
                log.Println("Stopping accept loop")
                wg.Wait()
                log.Printf("All connections closed, exiting")
                return
            default:
                log.Printf("Accept error: %v", err)
                continue
            }
        }
        // Handle connection
        wg.Add(1)
        go handleConnection(ctx, &wg, conn, store)
    }
}

// handle connection, parse request, log command, dispatch handler
func handleConnection(ctx context.Context, wg *sync.WaitGroup, conn net.Conn, store *datastore.Store) {
    defer wg.Done()
    defer conn.Close()

    reader := bufio.NewReader(conn)

    for {
        select {
        case <-ctx.Done():
            // Shutdown: stop processing requests
            log.Printf("Connection from %s closed due to shutdown", conn.RemoteAddr())
            return
        default:
            args, err := protocol.ParseRequest(reader)
            if err != nil {
                return
            }
            cmd := strings.ToUpper(string(args[0]))
            handler, ok := command.Handlers[cmd]
            var resp []byte
            if !ok {
                resp = protocol.FormatError(fmt.Sprintf("ERR unknown command '%s'", cmd))
            } else {
                resp = handler(store, args)
            }
            conn.Write(resp)
        }
    }
}