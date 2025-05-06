package main

import (
    "bufio"
    "fmt"
    "io"
    "log"
    "net"
    "strings"

    "github.com/CallumHomeDev/silk-kv/internal/command"
    "github.com/CallumHomeDev/silk-kv/internal/datastore"
    "github.com/CallumHomeDev/silk-kv/internal/protocol"
)

func main() {
    // Khởi store và listener
    store := datastore.New()
    listener, err := net.Listen("tcp", ":8081")
    if err != nil {
        log.Fatalf("Error starting server: %v", err)
    }
    defer listener.Close()
    log.Println("▶︎ Server listening on :8081")

    // Vòng accept kết nối
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("Error accepting connection: %v", err)
            continue
        }
        go handleConnection(conn, store)
    }
}

// Xử lý từng kết nối, parse request, log command và thực thi handler
func handleConnection(conn net.Conn, store *datastore.Store) {
    defer conn.Close()
    reader := bufio.NewReader(conn)

    for {
        args, err := protocol.ParseRequest(reader)
        if err != nil {
            if err != io.EOF {
                log.Printf("Error parsing request from %s: %v", conn.RemoteAddr(), err)
            }
            return
        }

        // Định danh command và args
        cmd := strings.ToUpper(string(args[0]))
        var argStrings []string
        for _, a := range args[1:] {
            argStrings = append(argStrings, string(a))
        }
        // In log ra server
        log.Printf("[%s] Command: %s, Args: %v", conn.RemoteAddr(), cmd, argStrings)

        // Dispatch handler
        handler, exists := command.Handlers[cmd]
        var resp []byte
        if !exists {
            resp = protocol.FormatError(fmt.Sprintf("ERR unknown command '%s'", cmd))
        } else {
            resp = handler(store, args)
        }
        // Gửi response
        conn.Write(resp)
    }
}