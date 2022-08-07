package main

import (
	"errors"
	"log"
	"net"
	"os"
	"time"

	"github.com/oklog/ulid/v2"
)

type client struct {
	id        ulid.ULID
	conn      net.Conn
	readChan  chan<- []byte
	writeChan chan []byte
}

func newClient(conn net.Conn, readChan chan<- []byte) *client {
	return &client{
		id:        ulid.Make(),
		conn:      conn,
		readChan:  readChan,
		writeChan: make(chan []byte, 100),
	}
}

func (c *client) readLoop() {
	content := make([]byte, 0, 1024)

	for {
		buf := make([]byte, 1024)
		c.conn.SetReadDeadline(time.Now().Add(50 * time.Millisecond))

		n, err := c.conn.Read(buf)
		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) {
				continue
			}

			log.Printf("error: (*client).readLoop: %s: %s", c.id, err)
			return
		}

		if n == 1024 {
			content = append(content, buf...)
			continue
		}

		content = append(content, buf[:n]...)
		c.readChan <- content

		log.Printf("(*client).readLoop: %s: %s", c.id, string(content))

		content = content[:0]
	}
}

func (c *client) write(content []byte) {
	c.writeChan <- content
}

func (c *client) writeLoop() {
	t := time.NewTicker(100 * time.Millisecond)

	for {
		select {
		case <-t.C:
			c.conn.SetWriteDeadline(time.Now().Add(50 * time.Millisecond))

			_, err := c.conn.Write([]byte("ping"))
			if err != nil {
				log.Printf("error: (*client).writeLoop: %s: %s", c.id, err)
				return
			}

		case content := <-c.writeChan:
			c.conn.SetWriteDeadline(time.Now().Add(100 * time.Millisecond))

			_, err := c.conn.Write(content)
			if err != nil {
				continue
			}

			log.Printf("(*client).writeLoop: %s: %s", c.id, string(content))
		}
	}
}
