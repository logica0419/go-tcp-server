package main

import (
	"log"
	"net"

	"github.com/oklog/ulid/v2"
)

type hub struct {
	clients  map[ulid.ULID]*client
	readChan chan []byte
}

func newHub() *hub {
	return &hub{
		clients:  make(map[ulid.ULID]*client),
		readChan: make(chan []byte, 100),
	}
}

func (h *hub) run() {
	for {
		content := <-h.readChan
		for _, c := range h.clients {
			c.write(content)
		}
	}
}

func (h *hub) connHandler(conn net.Conn) {
	c := newClient(conn, h.readChan)
	h.clients[c.id] = c

	log.Printf("client %s connected", c.id)

	go c.readLoop()
	c.writeLoop()

	conn.Close()
	delete(h.clients, c.id)

	log.Printf("client %s disconnected", c.id)
}
