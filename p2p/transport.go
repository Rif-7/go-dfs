package p2p

import "net"

// Represents a remote node.
type Peer interface {
	net.Conn
	Send([]byte) error
	CloseStream()
}

// Handles connections (TCP, UDP, Websockets..) between two nodes.
type Transport interface {
	Addr() string
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
	Dial(string) error
}
