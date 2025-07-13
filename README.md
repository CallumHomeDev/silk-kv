# silk-kv 🌿

`silk-kv` is a lightweight, in-memory key-value database written in Go.  
Its goal is to replicate core functionalities of **Redis**, including RESP protocol support, common data types, expiration features, and LRU-based memory management — while remaining simple and extensible.

This project is ideal for learning, prototyping, or building custom high-performance memory stores.

---

## 🚀 Features

- **Redis-compatible RESP protocol** – connect using `redis-cli` or compatible clients
- Key commands supported: `SET`, `GET`, `DEL`, `EXPIRE`, `TTL`, `INCR`, `DECR`, `LPUSH`, `RPOP`, etc.
- Single-threaded TCP server with non-blocking I/O (poll/epoll/kqueue)
- **LRU eviction** policy for automatic memory management
- Support for **key expiration (TTL)**
- Modular, extensible architecture with clean separation of logic and protocol
- Implemented in Go, using idiomatic and maintainable structure

---

## 📦 Getting Started

### Requirements

- Go 1.18+

### Installation

```bash
# 1. Clone the repository
git clone -b develop https://github.com/CallumHomeDev/silk-kv.git
cd silk-kv

# 2. Run the server
go run cmd/silk-kv/main.go

# 3. In another terminal, connect using redis-cli
redis-cli -p 8081

# 4. Try some commands
> PING
PONG

> SET hello world
OK

> GET hello
"world"
