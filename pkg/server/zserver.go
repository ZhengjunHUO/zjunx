package server

import (
	"net"
	"crypto/tls"
	"log"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ZhengjunHUO/zjunx/pkg/config"
	"github.com/ZhengjunHUO/zjunx/pkg/encoding"
)

type ZServer interface {
	Start()
	Stop()

	GetMux() ZMux
	GetCnxAdm() ZConnectionAdmin
	GetPostStartHook() func(ZConnection)
	GetPreStopHook() func(ZConnection)

	SetInterruptHandler()
}

type Server struct {
	Name		string
	ListenIP	string
	IPVersion	string
	ListenPort	uint64

	listener	net.Listener
	Mux		ZMux
	CnxAdm		ZConnectionAdmin

	PostStartHook	func(ZConnection)
	PreStopHook	func(ZConnection)
}

type ServerOption func(*Server)

func WithHandler(ct encoding.ZContentType, h ZHandler) ServerOption {
	return func(s *Server) {
		s.Mux.Register(ct, h)
	}
}

func WithPostStart(hook func(ZConnection)) ServerOption {
	return func(s *Server) {
		s.PostStartHook = hook
	}
}

func WithPreStop(hook func(ZConnection)) ServerOption {
	return func(s *Server) {
		s.PreStopHook = hook
	}
}

func NewServer(opts ...ServerOption) ZServer {
	s := &Server {
		Name: config.Cfg.ServerName,
		ListenIP: config.Cfg.ListenIP,
		IPVersion: "tcp4",
		ListenPort: config.Cfg.ListenPort,
		listener: nil,
		Mux: MuxInit(),
		CnxAdm: AdmInit(),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Lauch a ZJunx server
func (s *Server) Start() {
	log.Printf("[INFO] Starting %s ... \n", s.Name)

	// Load server's certificate and key
	cert, err := tls.LoadX509KeyPair("ssl/zjunx.crt", "ssl/zjunx.key")
	if err != nil {
		log.Fatalln(err)
	}

	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}

	// Parse the tcp addr from the server's config
	addr := fmt.Sprintf("%s:%d", s.ListenIP, s.ListenPort)
	if _, err = net.ResolveTCPAddr(s.IPVersion, addr); err != nil {
		log.Fatalln("[FATAL] ", err)
	}

	// Bring up a tcp network listener with TLS enabled
	s.listener, err = tls.Listen("tcp", addr, cfg)
	if err != nil {
		log.Fatalln("[FATAL] ", err)
	}
	defer s.listener.Close()
	log.Printf("[INFO] Server is up, listening at %s:%d, max connection: %v\n", s.ListenIP, s.ListenPort, config.Cfg.ConnLimit)
	log.Printf("[DEBUG] Workers: %v; Queue length per worker: %v; Scheduling algorithm: %v", config.Cfg.WorkerProcesses, config.Cfg.BacklogSize, config.Cfg.ScheduleAlgo)

	go s.SetInterruptHandler()
	var cnxID uint64
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			//log.Println("[WARN] ", err)
			continue
		}

		log.Println("[DEBUG] Current connection number: ", s.CnxAdm.PoolSize())
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
	s.listener.Close()
	s.CnxAdm.Evacuate()
	s.Mux.WorkerDismiss()
	log.Printf("[INFO] %s stopped.\n", s.Name)
	os.Exit(0)
}

func (s *Server) SetInterruptHandler() {
	chIntr := make(chan os.Signal, 1)
	signal.Notify(chIntr, os.Interrupt, syscall.SIGTERM)

	<-chIntr
	log.Println("Interrupt singal catched!")
	s.Stop()
}

func (s *Server) GetPostStartHook() func(ZConnection) {
	return s.PostStartHook
}

func (s *Server) GetPreStopHook() func(ZConnection) {
	return s.PreStopHook
}
