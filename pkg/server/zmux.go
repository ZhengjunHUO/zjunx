package server

import (
	"log"

	"github.com/ZhengjunHUO/zjunx/pkg/encoding"
	"github.com/ZhengjunHUO/zjunx/pkg/config"
)

type ZMux interface {
	WorkerInit()
	Register(encoding.ZContentType, ZHandler)
	Schedule(ZRequest)	
	Handle(ZRequest)
}

type Mux struct {
	WorkerProcesses uint64
	WorkerBacklog	[]chan ZRequest
	HandlerSet 	map[encoding.ZContentType]ZHandler
	ScheduleAlgo	string
}

func MuxInit() ZMux {
	return &Mux{
		WorkerProcesses: config.Cfg.WorkerProcesses,
		WorkerBacklog: make([]chan ZRequest, config.Cfg.WorkerProcesses),
		HandlerSet: make(map[encoding.ZContentType]ZHandler),
		ScheduleAlgo: config.Cfg.ScheduleAlgo, 
	}
}

func (m *Mux) WorkerInit() {
	for i := range m.WorkerBacklog {
		m.WorkerBacklog[i] = make(chan ZRequest, config.Cfg.BacklogSize)
		go func(wid int, backlog chan ZRequest){
			log.Printf("Worker %d up.\n", wid)
			for {
				select {
					case req := <-backlog :
					m.Handle(req)
				}
			}			
		}(i, m.WorkerBacklog[i])	
	}
}

func (m *Mux) Register(ct encoding.ZContentType, h ZHandler) {
	if _, ok := m.HandlerSet[ct]; ok {
		log.Printf("Handler %d found, will be overwritten.\n")
	}

	m.HandlerSet[ct] = h
	log.Printf("Handler %d registered.\n")
}

func (m *Mux) Schedule(req ZRequest) {
	// TO IMPLEMENT
	// switch m.ScheduleAlgo
	m.WorkerBacklog[0] <- req	
}

func (m *Mux) Handle(req ZRequest) {
	h, ok := m.HandlerSet[req.GetContentType()]
	if ok {
		h.Handle(req)
	}else{
		log.Printf("Unknown content type (%d) from request, skip.\n", req.GetContentType())
	}
}
