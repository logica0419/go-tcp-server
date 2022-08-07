package main

import (
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
		select {
		case content := <-h.readChan:
			for _, c := range h.clients {
				c.write(content)
			}
		}
	}

}

func (h *hub) runClient(conn net.Conn) {
	c := newClient(conn, h.readChan)
	h.clients[c.id] = c

	go c.readLoop()
	c.writeLoop()

	conn.Close()
	delete(h.clients, c.id)
}
