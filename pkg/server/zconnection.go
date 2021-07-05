package server

import (
	"net"
	"io"
	"log"
	"errors"

	"github.com/ZhengjunHUO/zjunx/pkg/encoding"
)

type ZConnection interface {
	Start()
	Reader()
	Writer()
	Close()
	RespondToClient(encoding.ZContentType, []byte) error 
	GetID() uint64
}

type Connection struct {
	ID	uint64
	Conn	*net.TCPConn	
	Server	ZServer
	Mux	ZMux

	chServerResp	chan []byte
	chClose		chan bool
	isActive	bool
}

func ConnInit(cnxID uint64, conn *net.TCPConn, s ZServer, m ZMux) *Connection {
	return &Connection{
		ID: cnxID,
		Conn: conn,
		Server: s,
		Mux: m,
		chServerResp: make(chan []byte),
		chClose: make(chan bool, 1),
		isActive: true,
	}
}

// Read from the TCP stream payload and decode the raw bytes to struct
// Prepare a processed request and send it to a worker to handle it
func (c *Connection) Reader() {
	defer c.Close()

	blk := encoding.BlockInit()
	for {
		ct := encoding.ContentInit(encoding.ZContentType(0), []byte{})
		if err := blk.Unmarshalling(c.Conn, ct); err != nil {
			if err != io.EOF {
				log.Println("[WARN] Unmarshalling failed: ", err)
			}
			break
		}

		req := ReqInit(c, ct)
		c.Mux.Schedule(req)
		log.Println("[DEBUG] Request sheduled.")
	}

}

// Called by handler after dealing with the request
// Send the raw bytes to Writer
func (c *Connection) RespondToClient(ct encoding.ZContentType, data []byte) error {
	if ! c.isActive	{
		return errors.New("Error sending response: Connection is closed")
	}

	blk := encoding.BlockInit()
	buf, err := blk.Marshalling(encoding.ContentInit(ct, data))
	if err != nil {
		return errors.New("Error sending response to client !")	
	}

	c.chServerResp <- buf
	return nil
}

// Write the processed response (raw bytes received from handler) to client
// Quit on receiving the close signal after Reader quits
func (c *Connection) Writer() {
	for {
		select {
			case data := <- c.chServerResp:
			if _, err := c.Conn.Write(data); err != nil {
				log.Println("[ERROR] Write back to client: ", err)
				return
			}
			case <- c.chClose:
				return
		}
	}
}

// Seperate read/write thread, leave the handling part to ZMux
func (c *Connection) Start() {
	log.Printf("[DEBUG] Connection [id: %d] established from %v\n", c.ID, c.Conn.RemoteAddr())
	go c.Reader()
	go c.Writer()
}

func (c *Connection) GetID() uint64 {
	return c.ID
}

// Clenup current connection before exit
func (c *Connection) Close() {
	if !c.isActive {
		return
	}
	c.isActive = false

	c.chClose <- true
	close(c.chClose)
	close(c.chServerResp)

	if err := c.Conn.Close(); err !=nil {
		log.Printf("[DEBUG] Closing connection [id: %d]: %s", c.ID, err)
	}else{
		log.Printf("[DEBUG] Connection [id: %d] closed. \n", c.ID)
	}
}
