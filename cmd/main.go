package main

import (
	"flag"
	"log"
	"net"

	"github.com/papertigers/tcpproxy"
)

var (
	localAddrPtr  = flag.String("l", ":8080", "local address")
	remoteAddrPtr = flag.String("r", "localhost:80", "remote address")
)

func main() {
	flag.Parse()

	laddr, err := net.ResolveTCPAddr("tcp", *localAddrPtr)
	if err != nil {
		log.Fatalf("Failed to resolve local address: %s\n", err)
	}

	raddr, err := net.ResolveTCPAddr("tcp", *remoteAddrPtr)
	if err != nil {
		log.Fatalf("Failed to resolve remote address: %s\n", err)
	}

	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		log.Fatalf("Failed to open local port: %s\n", err)
	}

	// Setup a proxy for each incomming connection
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("Failed to accept connection: %s\n", err)
			continue
		}

		p := proxy.New(conn, laddr, raddr)
		go p.Start()
	}
}
