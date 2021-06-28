package server

type ZConnectionAdmin interface {
	Register(ZConnection)
	Retrieve(uint64) ZConnection
	Remove(ZConnection)
	Evacuate()
	PoolSize() int
}
