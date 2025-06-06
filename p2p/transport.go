package p2p

// Represents a remote node.
type Peer interface {
	Close() error
}

// Handles connections (TCP, UDP, Websockets..) between two nodes.
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
