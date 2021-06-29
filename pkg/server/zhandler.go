package server

type ZHandler interface {
	Handle(ZRequest)
}

type Handler struct {}

func (h *Handler) Handle(req ZRequest) {}
