package server

import (
	"sync"
)

type ZConnectionAdmin interface {
	Register(ZConnection)
	Retrieve(uint64) ZConnection
	Remove(ZConnection)
	Evacuate()
	PoolSize() int
}

type ConnectionAdmin struct {
	Pool	map[uint64]ZConnection
	Mutex	sync.RWMutex
}

func AdmInit() ZConnectionAdmin {
	return &ConnectionAdmin {
		Pool: make(map[uint64]ZConnection),
	}
}

func (ca *ConnectionAdmin) Register(conn ZConnection) {
	ca.Mutex.Lock()
	defer ca.Mutex.Unlock()

	ca.Pool[conn.GetID()] = conn
}

func (ca *ConnectionAdmin) Retrieve(cid uint64) ZConnection {
	ca.Mutex.RLock()
	defer ca.Mutex.RUnlock()

	if conn, ok := ca.Pool[cid]; ok {
		return conn
	} else {
		return nil
	}
}

func (ca *ConnectionAdmin) Remove(conn ZConnection) {
	ca.Mutex.Lock()
	defer ca.Mutex.Unlock()

	delete(ca.Pool, conn.GetID())
}

func (ca *ConnectionAdmin) Evacuate() {
	ca.Mutex.Lock()
	defer ca.Mutex.Unlock()

	for cid, conn := range ca.Pool {
		conn.Close()
		delete(ca.Pool, cid)
	}
}

func (ca *ConnectionAdmin) PoolSize() int {
	return len(ca.Pool)
}
