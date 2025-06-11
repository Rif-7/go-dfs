package main

import (
	"bytes"
	"fmt"
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

	for i := range 10 {
		data := bytes.NewReader([]byte("your private data here"))
		if err := s2.Store(fmt.Sprintf("yourkeyhere_%d", i), data); err != nil {
			fmt.Println(err)
		}
		time.Sleep(5 * time.Millisecond)
	}

	// r, err := s2.Get("yourkeyhere")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// b, err := io.ReadAll(r)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(string(b))

	select {}
}
