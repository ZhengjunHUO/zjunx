package server

import (
	"net"
	"log"
	"fmt"

	"github.com/ZhengjunHUO/zjunx/pkg/config"
)

type ZServer interface {
	Start()
	Stop()
}

type Server struct {
	Name		string
	ListenIP	string
	IPVersion	string
	ListenPort	uint16
	
	Mux		ZMux
}

func ServerInit() ZServer {
	return &Server {
		Name: config.Cfg.ServerName,
		ListenIP: config.Cfg.ListenIP,
		IPVersion: "tcp4",
		ListenPort: config.Cfg.ListenPort,
		Mux: MuxInit(),
	}
}

func (s *Server) Start() {
	log.Printf("[INFO] Starting %s ... \n", s.Name)

	// Parse the tcp addr from the server's config
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.ListenIP, s.ListenPort))
	if err != nil {
		log.Fatalln("[FATAL] ", err)
	}

	// launch a tcp network listener
	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		log.Fatalln("[FATAL] ", err)
	}
	defer listener.Close()
	log.Printf("[INFO] Server is up, listening at %s:%d\n", s.ListenIP, s.ListenPort)

	var cnxID uint64
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("[WARN] ", err)
			continue
		}

		cnx := ConnInit(cnxID, conn, s)
		go cnx.Start()
		cnxID += 1
	}
}

func (s *Server) Stop() {
	log.Printf("[INFO] %s stopped.\n", s.Name)
}
