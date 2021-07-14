package server

import (
	"sync"
	"log"
)

type ZConnectionAdmin interface {
	Register(ZConnection)
	Retrieve(uint64) ZConnection
	RetrieveAll() []ZConnection
	Remove(ZConnection)
	Evacuate()
	PoolSize() int
}

type ConnectionAdmin struct {
	// A vault of connections
	Pool	map[uint64]ZConnection
	Mutex	sync.RWMutex
}

func AdmInit() ZConnectionAdmin {
	return &ConnectionAdmin {
		Pool: make(map[uint64]ZConnection),
	}
}

// Add new connection to the vault
func (ca *ConnectionAdmin) Register(conn ZConnection) {
	ca.Mutex.Lock()
	defer ca.Mutex.Unlock()

	ca.Pool[conn.GetID()] = conn
}

// Get the connection via its id
func (ca *ConnectionAdmin) Retrieve(cid uint64) ZConnection {
	ca.Mutex.RLock()
	defer ca.Mutex.RUnlock()

	if conn, ok := ca.Pool[cid]; ok {
		return conn
	} else {
		return nil
	}
}

func (ca *ConnectionAdmin) RetrieveAll() []ZConnection {
	ca.Mutex.RLock()
	defer ca.Mutex.RUnlock()

	conns := make([]ZConnection, 0, len(ca.Pool))

	for _, v := range ca.Pool {
		conns = append(conns, v)
	}

	return conns
}

// Delete the connection from the vault
func (ca *ConnectionAdmin) Remove(conn ZConnection) {
	ca.Mutex.Lock()
	defer ca.Mutex.Unlock()

	delete(ca.Pool, conn.GetID())
}

// Called when daemon is shut down, closing all active connections to ensure a clean quit
func (ca *ConnectionAdmin) Evacuate() {
	ca.Mutex.Lock()
	defer ca.Mutex.Unlock()
	log.Println("[INFO] Stop all connections ...")

	for cid, conn := range ca.Pool {
		conn.Close()
		delete(ca.Pool, cid)
	}
}

func (ca *ConnectionAdmin) PoolSize() int {
	return len(ca.Pool)
}
