package server

import (
	"log"

	"github.com/ZhengjunHUO/zjunx/pkg/encoding"
	"github.com/ZhengjunHUO/zjunx/pkg/config"
)

// A multiplexer dealing with client's requests
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

// Initialize a pool of worker processes to handle requests
// Each worker process is assigned with a buffer queue
func (m *Mux) WorkerInit() {
	for i := range m.WorkerBacklog {
		m.WorkerBacklog[i] = make(chan ZRequest, config.Cfg.BacklogSize)
		go func(wid int, backlog chan ZRequest){
			log.Printf("[DEBUG] Worker %d up.\n", wid)
			for {
				select {
					case req := <-backlog :
					m.Handle(req)
				}
			}			
		}(i, m.WorkerBacklog[i])	
	}
}

// Register the serverside defined handler to ZMux
func (m *Mux) Register(ct encoding.ZContentType, h ZHandler) {
	if _, ok := m.HandlerSet[ct]; ok {
		log.Printf("[DEBUG] Handler %v found, will be overwritten.\n", ct)
	}

	m.HandlerSet[ct] = h
	log.Printf("[DEBUG] Handler %v registered.\n", ct)
}

// Scheduling the request to appropriate worker depending on the algorithm
func (m *Mux) Schedule(req ZRequest) {
	// TO IMPLEMENT
	// switch m.ScheduleAlgo
	m.WorkerBacklog[0] <- req	
}

// Handle client's request if related handler is registered
func (m *Mux) Handle(req ZRequest) {
	h, ok := m.HandlerSet[req.ContentType()]
	if ok {
		h.Handle(req)
	}else{
		log.Printf("[WARN] Unknown content type (%d) from request, skip.\n", req.ContentType())
	}
}
