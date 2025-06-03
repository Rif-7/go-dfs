package p2p

import "net"

// Holds data being being sent between two nodes in the network
type RPC struct {
	From    net.Addr
	Payload []byte
}
