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
	s3 := makeServer(":5000", "5000", ":4000", ":3000")

	go func() { log.Fatal(s1.Start()) }()
	time.Sleep(time.Second)
	go func() { log.Fatal(s2.Start()) }()
	time.Sleep(time.Second)
	go func() { log.Fatal(s3.Start()) }()
	time.Sleep(time.Second)

	for i := range 5 {
		key := fmt.Sprintf("helloworld_%d", i)
		data := bytes.NewReader([]byte("your private data here"))
		if err := s3.Store(key, data); err != nil {
			log.Fatal(err)
		}

		if err := s3.store.Delete(s3.ID, key); err != nil {
			log.Fatal(err)
		}

		r, err := s3.Get(key)
		if err != nil {
			log.Fatal(err)
		}

		b, err := io.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}

		if rc, ok := r.(io.ReadCloser); ok {
			rc.Close()
		}

		fmt.Println("data: ", string(b))

		// Comment this condition out if you want
		// to see the contents of the files
		if err := s3.Delete(key); err != nil {
			log.Fatal(err)
		}

	}

	select {}

}
