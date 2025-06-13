package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/Rif-7/go-dfs/p2p"
)

func makeServer(listenAddr string, root string, nodes ...string) *FileServer {
	tcptransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}

	tcpTransport := p2p.NewTCPTransport(tcptransportOpts)

	fileServerOpts := FileServerOpts{
		EncKey:            newEncryptionKey(),
		StorageRoot:       root + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}

	s := NewFileServer(fileServerOpts)
	tcpTransport.OnPeer = s.OnPeer

	return s

}

func main() {
	s1 := makeServer(":3000", "3000", "")
	s2 := makeServer(":4000", "4000", ":3000")

	go func() {
		log.Fatal(s1.Start())
	}()

	time.Sleep(time.Second)

	go s2.Start()
	time.Sleep(time.Second)

	key := "helloworld"
	data := bytes.NewReader([]byte("your private data here"))
	if err := s2.Store(key, data); err != nil {
		fmt.Println(err)
	}

	if err := s2.store.Delete(key); err != nil {
		log.Fatal(err)
	}

	r, err := s2.Get(key)
	if err != nil {
		log.Fatal(err)
	}

	b, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))

	select {}
}
