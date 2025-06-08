package p2p

import "net"

// Represents a remote node.
type Peer interface {
	Send([]byte) error
	RemoteAddr() net.Addr
	Close() error
}

// Handles connections (TCP, UDP, Websockets..) between two nodes.
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
	Dial(string) error
}
