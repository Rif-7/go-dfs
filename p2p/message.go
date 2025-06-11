package p2p

const (
	IncomingMessage = 0x1
	IncomingStream  = 0x2
)

// Holds data being being sent between two nodes in the network
type RPC struct {
	From    string
	Payload []byte
	Stream  bool
}
