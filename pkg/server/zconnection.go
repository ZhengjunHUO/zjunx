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
	GetServer() ZServer
	UpdateContext(string, interface{})
	GetContext(string) interface{}
	DeleteContext(string)
}

type Connection struct {
	ID	uint64
	Conn	*net.TCPConn	
	Server	ZServer
	Context map[string]interface{}

	chServerResp	chan []byte
	chClose		chan bool
	isActive	bool
}

func ConnInit(cnxID uint64, conn *net.TCPConn, s ZServer) ZConnection {
	cnx := &Connection{
		ID: cnxID,
		Conn: conn,
		Server: s,
		Context: make(map[string]interface{}),
		chServerResp: make(chan []byte),
		chClose: make(chan bool, 1),
		isActive: true,
	}

	s.GetCnxAdm().Register(cnx)
	return cnx
}

// Read from the TCP stream payload and decode the raw bytes to struct
// Prepare a processed request and send it to a worker to handle it
func (c *Connection) Reader() {
	defer c.Close()
	defer c.Server.GetCnxAdm().Remove(c)

	blk := encoding.BlockInit()
	for {
		ct := encoding.ContentInit(encoding.ZContentType(0), []byte{})
		if err := blk.Unmarshalling(c.Conn, ct); err != nil {
			if err != io.EOF {
				log.Println("[WARN] Unmarshalling failed: ", err)
			}
			break
		}

		// build a request from a valid incoming package
		req := ReqInit(c, ct)
		// sending request to a worker
		c.Server.GetMux().Schedule(req)
		log.Printf("[DEBUG] Request from Connection [id: %v] sheduled.\n", c.GetID())
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

	// Writer process listening at the other end of channel
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
	c.Server.CallPostStart(c)
}

func (c *Connection) GetID() uint64 {
	return c.ID
}

func (c *Connection) GetServer() ZServer {
	return c.Server
}

func (c *Connection) UpdateContext(key string, value interface{}) {
	c.Context[key] = value
}

func (c *Connection) GetContext(key string) interface{} {
	if v, ok := c.Context[key]; ok {
		return v
	}else{
		return nil
	}
}

func (c *Connection) DeleteContext(key string) {
	delete(c.Context, key)
}

// Clenup current connection before exit
func (c *Connection) Close() {
	log.Printf("[DEBUG] Closing connection [id: %d] ... \n", c.ID)
	if !c.isActive {
		return
	}
	c.isActive = false

	// close the Writer process
	c.chClose <- true
	close(c.chClose)
	close(c.chServerResp)

	if err := c.Conn.Close(); err !=nil {
		log.Printf("[DEBUG] Closing connection [id: %d]: %s\n", c.ID, err)
	}else{
		log.Printf("[DEBUG] Connection [id: %d] closed. \n", c.ID)
	}
}
