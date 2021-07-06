package server

import (
	"net"
	"log"
	"fmt"

	"github.com/ZhengjunHUO/zjunx/pkg/config"
	"github.com/ZhengjunHUO/zjunx/pkg/encoding"
)

type ZServer interface {
	Start()
	Stop()
	RegistHandler(encoding.ZContentType, ZHandler)
	GetMux() ZMux
	GetCnxAdm() ZConnectionAdmin
}

type Server struct {
	Name		string
	ListenIP	string
	IPVersion	string
	ListenPort	uint16
	
	Mux		ZMux
	CnxAdm		ZConnectionAdmin
}

func ServerInit() ZServer {
	return &Server {
		Name: config.Cfg.ServerName,
		ListenIP: config.Cfg.ListenIP,
		IPVersion: "tcp4",
		ListenPort: config.Cfg.ListenPort,
		Mux: MuxInit(),
		CnxAdm: AdmInit(),
	}
}

// Lauch a ZJunx server
func (s *Server) Start() {
	log.Printf("[INFO] Starting %s ... \n", s.Name)

	// Parse the tcp addr from the server's config
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.ListenIP, s.ListenPort))
	if err != nil {
		log.Fatalln("[FATAL] ", err)
	}

	// Bring up a tcp network listener
	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		log.Fatalln("[FATAL] ", err)
	}
	defer listener.Close()
	log.Printf("[INFO] Server is up, listening at %s:%d\n", s.ListenIP, s.ListenPort)

	// Initialize a working pool to handler requests
	s.Mux.WorkerInit()

	var cnxID uint64
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("[WARN] ", err)
			continue
		}

		if uint64(s.CnxAdm.PoolSize()) < config.Cfg.ConnLimit {
			// Create a Goroutine for each incoming client
			cnx := ConnInit(cnxID, conn, s)
			go cnx.Start()
			cnxID += 1
		}else{
			log.Println("[INFO] Connection number reach limit.")
			blk := encoding.BlockInit()
			buf, _ := blk.Marshalling(encoding.ContentInit(encoding.ZContentType(0), []byte("Server busy, please try again later ...")))
			conn.Write(buf)
			conn.Close()
		}
	}
}

func (s *Server) GetMux() ZMux {
        return s.Mux
}

func (s *Server) GetCnxAdm() ZConnectionAdmin {
	return s.CnxAdm
}

func (s *Server) Stop() {
	s.CnxAdm.Evacuate()
	log.Printf("[INFO] %s stopped.\n", s.Name)
}

func (s *Server) RegistHandler(ct encoding.ZContentType, h ZHandler) {
	s.Mux.Register(ct, h)
}
