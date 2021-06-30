package server

import (
	"net"
	"fmt"
	"log"
	"io"
	"os"
)

type ZConnection interface {
	Start()
	Reader()
	Writer()
	Close()
}

type Connection struct {
	ID	uint64
	Conn	*net.TCPConn	
	Server	ZServer
}

func ConnInit(cnxID uint64, conn *net.TCPConn, s ZServer) *Connection {
	return &Connection{
		ID: cnxID,
		Conn: conn,
		Server: s,
	}
}

func (c *Connection) Reader() {
	defer c.Close()
	for {
		if _, err := io.Copy(os.Stdout, c.Conn); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func (c *Connection) Writer() {

}

func (c *Connection) Start() {
	log.Printf("[DEBUG] Connection [id: %d] established from %v\n", c.ID, c.Conn.RemoteAddr())
	go c.Reader()
//	go c.Writer()
}

func (c *Connection) Close() {
	if err := c.Conn.Close(); err !=nil {
		log.Printf("[DEBUG] Closing connection [id: %d]: %s", c.ID, err)
	}else{
		log.Printf("[DEBUG] Connection [id: %d] closed. \n", c.ID)
	}
}
