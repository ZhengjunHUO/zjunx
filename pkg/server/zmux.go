package server

import (
	"log"
	"time"
	"math/rand"

	"github.com/ZhengjunHUO/zjunx/pkg/encoding"
	"github.com/ZhengjunHUO/zjunx/pkg/config"
)

// A multiplexer dealing with client's requests
type ZMux interface {
	Register(encoding.ZContentType, ZHandler)
	Schedule(ZRequest)	
	Handle(ZRequest)
	WorkerDismiss()
}

type Mux struct {
	WorkerProcesses uint64
	WorkerBacklog	[]chan ZRequest
	WorkerExit	[]chan bool
	HandlerSet 	map[encoding.ZContentType]ZHandler
	// legitime value: RoundRobin, Random, LeastConn
	ScheduleAlgo	string
}

func MuxInit() ZMux {
	m := &Mux{
		WorkerProcesses: config.Cfg.WorkerProcesses,
		WorkerBacklog: make([]chan ZRequest, config.Cfg.WorkerProcesses),
		WorkerExit: make([]chan bool, config.Cfg.WorkerProcesses),
		HandlerSet: make(map[encoding.ZContentType]ZHandler),
		ScheduleAlgo: config.Cfg.ScheduleAlgo, 
	}

	// Initialize a pool of worker processes to handle requests
	// Each worker process is assigned with a buffer queue

	for i := range m.WorkerBacklog {
		m.WorkerBacklog[i] = make(chan ZRequest, config.Cfg.BacklogSize)
		m.WorkerExit[i] = make(chan bool)
		go func(wid int, backlog chan ZRequest, chExit chan bool){
			log.Printf("[DEBUG] Worker %d up.\n", wid)
			mainloop: for {
				select {
					case req := <-backlog :
						m.Handle(req)
					case <-chExit:
						break mainloop
				}
			}			
			log.Printf("[DEBUG] Worker %d dismissed.\n", wid)
		}(i, m.WorkerBacklog[i], m.WorkerExit[i])
	}

	return m
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
	switch m.ScheduleAlgo {
		case "Random":
			m.WorkerBacklog[rand.Intn(int(m.WorkerProcesses))] <- req
			/* Debug
			temp := rand.Intn(int(m.WorkerProcesses))
			log.Println("Send to : ", temp)
			m.WorkerBacklog[temp] <- req
			*/
		case "LeastConn":
			min, target := int(config.Cfg.BacklogSize), 0

			for i := range m.WorkerBacklog {
				if qs := len(m.WorkerBacklog[i]); qs == 0 {
					//log.Println("Send to : ", i)
					m.WorkerBacklog[i] <- req
					return
				}else{
					if qs < min {
						min, target = qs, i
					}
				}
			}

			//log.Println("Send to : ", target)
			m.WorkerBacklog[target] <- req
		default:
			// default using "RoundRobin"
			//log.Println("Send to : ", req.Connection().GetID()%m.WorkerProcesses)
			m.WorkerBacklog[req.Connection().GetID()%m.WorkerProcesses] <- req
	}
}

// Handle client's request if related handler is registered
func (m *Mux) Handle(req ZRequest) {
	h, ok := m.HandlerSet[req.ContentType()]
	if ok {
		h.Handle(req)
	}else{
		log.Printf("[WARN] Unknown content type (%d) from request, skip.\n", req.ContentType())
		req.Connection().RespondToClient(encoding.ZContentType(0), []byte("Unknown request !\n"))
	}
}

func (m *Mux) WorkerDismiss() {
	for i := range m.WorkerExit {
		m.WorkerExit[i] <- true
		time.Sleep(time.Millisecond * 500)
		close(m.WorkerExit[i])
		close(m.WorkerBacklog[i])
	}

	log.Printf("[DEBUG] All workers are dismissed.\n")
}
