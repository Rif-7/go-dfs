# Go-DFS: A Lightweight Distributed File System in Go

**Go-DFS** is a peer-to-peer distributed file system written in Go. Designed for clarity, security, and modularity, it enables multiple nodes to share and replicate files across a network with encrypted streams and plug-and-play transport logic.

---

## Features

- Modular Transport Layer: Built-in TCP with support for adding UDP/WebSockets  
- AES-CTR Encrypted Streams: Secure file transfers across nodes  
- Bootstrapped Peer Discovery: Dynamically connect to known peers  
- Content-Addressed Storage: Efficient storage layout based on SHA1/MD5  
- Event-Driven Messaging: Simple RPC-based messaging framework  
- Pluggable Decoding: Choose between GOB and raw format decoding  
- Unit Tested: Coverage for core components including storage, crypto, and transport  

---

## Project Structure

```
go-dfs/
│
├── main.go               # Entry point - spins up multiple file servers
├── server.go             # Core FileServer logic and peer communication
├── store.go              # CAS disk storage engine
├── crypto.go             # AES encryption & decryption helpers
│
├── p2p/                  # Peer-to-peer networking layer
│   ├── transport.go           # Interfaces for Peer and Transport
│   ├── tcp_transport.go       # TCP-based transport implementation
│   ├── encoding.go            # Decoder abstraction (GOB/raw)
│   └── handshake.go           # Handshake logic
│
├── *.test.go             # Unit tests for each module
├── Makefile              # Build, Run, and Test helper
└── README.md             # You're reading it
```

---

## How It Works

- `main.go` sets up and runs multiple `FileServer` instances, simulating a distributed environment.
- Each `FileServer` manages:
  - Local file storage with `store.go` using a content-addressed layout.
  - Peer communication via the `p2p` package (`TCPTransport` and `Peer` interfaces).
  - File sharing using RPC-style messages: store, retrieve, and stream files across peers.
- Files are encrypted with AES-CTR before being sent over the network (`crypto.go`).
- `FileServer.Store()` saves and propagates files to peers.
- `FileServer.Get()` retrieves the file either locally or from peers, decrypting it before use.
- `FileServer.Delete()` deletes the file from local disk and from peers.

---

## Requirements

- Go 1.20+ (Tested on Go 1.22.5)  
- No external dependencies — pure standard library

---

## Getting Started

Clone the repository.

### Build the Binary

```bash
make build
```

### Run the Demo

```bash
make run
```

This launches three file servers on ports `:3000`, `:4000`, and `:5000`. Files are added on `:5000`, deleted, then re-fetched from the network via peer discovery and encrypted streams.

### Run Tests

```bash
make test
```

Runs unit tests for all supported modules.

---

## Future Improvements

- Peer discovery via DHT or multicast  
- File replication enforcement  
