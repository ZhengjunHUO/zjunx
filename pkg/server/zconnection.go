package server

import (
	"net"
	"log"

	"github.com/ZhengjunHUO/zjunx/pkg/encoding"
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
	Mux	ZMux
}

func ConnInit(cnxID uint64, conn *net.TCPConn, s ZServer, m ZMux) *Connection {
	return &Connection{
		ID: cnxID,
		Conn: conn,
		Server: s,
		Mux: m,
	}
}

func (c *Connection) Reader() {
	defer c.Close()

	blk := encoding.BlockInit()
	for {
		ct := encoding.ContentInit(encoding.ZContentType(0), []byte{})
		if err := blk.Unmarshalling(c.Conn, ct); err != nil {
			log.Println("[WARN] Unmarshalling failed: ", err)
			break
		}

		req := ReqInit(c, ct)
		c.Mux.Schedule(req)
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
