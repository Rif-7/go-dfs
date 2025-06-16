package p2p

import (
	"testing"
)

func TestTCPTransport(t *testing.T) {
	listenAddr := ":3000"

	opts := TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: NOPHandshakeFunc,
		Decoder:       DefaultDecoder{},
	}
	tr := NewTCPTransport(opts)

	if tr.ListenAddr != listenAddr {
		t.Errorf("expected ListenAddr to be %s, got %s", listenAddr, tr.ListenAddr)
	}

	if err := tr.ListenAndAccept(); err != nil {
		t.Errorf("expected ListenAndAccept to return nil, got %s", err)
	}

}
