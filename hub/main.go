package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	port := 8000
	protocol := "tcp"

	addr, err := net.ResolveTCPAddr(protocol, fmt.Sprintf(":%d", port))
	if err != nil {
		log.Panic(err)
	}

	sock, err := net.ListenTCP(protocol, addr)
	if err != nil {
		log.Panic(err)
	}
	defer sock.Close()

	log.Printf("Listening on %s\n", sock.Addr())

	h := newHub()
	go h.run()

	for {
		conn, err := sock.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go h.runClient(conn)
	}
}
