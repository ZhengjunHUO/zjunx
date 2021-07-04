package server

type ZHandler interface {
	Handle(ZRequest)
}

// Define the behavior of server in response to client's request
// Need to be implemented at server side 
// and registerd to ZJunx server's multiplexer (ZMux) after server is up
type Handler struct {}

func (h *Handler) Handle(req ZRequest) {}
