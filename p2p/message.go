package p2p

// Holds data being being sent between two nodes in the network
type RPC struct {
	From    string
	Payload []byte
}
