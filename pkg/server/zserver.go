package server

import (
	"net"
	"log"
	"fmt"
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
}

func ServerInit() ZServer {
	return &Server {
		Name: "ZJunx Server",
		ListenIP: "127.0.0.1",
		IPVersion: "tcp4",
		ListenPort: 8080,
	}
}

func (s *Server) Start() {
	log.Printf("Starting %s ... \n", s.Name)

	// Parse the tcp addr from the server's config
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.ListenIP, s.ListenPort))
	if err != nil {
		log.Fatalln(err)
	}

	// launch a tcp network listener
	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()
	log.Printf("Server is up, listening at %s:%d\n", s.ListenIP, s.ListenPort)

	var cnxID uint64
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}

		cnx := ConnInit(cnxID, conn, s)
		go cnx.Start()
		cnxID += 1
	}
}

func (s *Server) Stop() {
	log.Printf("%s stopped.\n", s.Name)
}
